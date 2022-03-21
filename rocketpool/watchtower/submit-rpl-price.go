package watchtower

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"github.com/rocket-pool/rocketpool-go/utils/client"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/contracts"
	rpgas "github.com/rocket-pool/smartnode/shared/services/gas"
	"github.com/rocket-pool/smartnode/shared/services/wallet"
	"github.com/rocket-pool/smartnode/shared/utils/api"
	"github.com/rocket-pool/smartnode/shared/utils/log"
	mathutils "github.com/rocket-pool/smartnode/shared/utils/math"
)

// Settings
const SubmitFollowDistancePrices = 2
const ConfirmDistancePrices = 30

// Submit GGP price task
type submitGgpPrice struct {
	c              *cli.Context
	log            log.ColorLogger
	cfg            config.RocketPoolConfig
	ec             *client.EthClientProxy
	w              *wallet.Wallet
	rp             *rocketpool.RocketPool
	oio            *contracts.OneInchOracle
	maxFee         *big.Int
	maxPriorityFee *big.Int
	gasLimit       uint64
}

// Create submit GGP price task
func newSubmitGgpPrice(c *cli.Context, logger log.ColorLogger) (*submitGgpPrice, error) {

	// Get services
	cfg, err := services.GetConfig(c)
	if err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	ec, err := services.GetEthClientProxy(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}
	oio, err := services.GetOneInchOracle(c)
	if err != nil {
		return nil, err
	}

	// Get the user-requested max fee
	maxFee, err := cfg.GetMaxFee()
	if err != nil {
		return nil, fmt.Errorf("Error getting max fee in configuration: %w", err)
	}

	// Get the user-requested max fee
	maxPriorityFee, err := cfg.GetMaxPriorityFee()
	if err != nil {
		return nil, fmt.Errorf("Error getting max priority fee in configuration: %w", err)
	}
	if maxPriorityFee == nil || maxPriorityFee.Uint64() == 0 {
		logger.Println("WARNING: priority fee was missing or 0, setting a default of 2.")
		maxPriorityFee = big.NewInt(2)
	}

	// Get the user-requested gas limit
	gasLimit, err := cfg.GetGasLimit()
	if err != nil {
		return nil, fmt.Errorf("Error getting gas limit in configuration: %w", err)
	}

	// Return task
	return &submitGgpPrice{
		c:              c,
		log:            logger,
		cfg:            cfg,
		ec:             ec,
		w:              w,
		rp:             rp,
		oio:            oio,
		maxFee:         maxFee,
		maxPriorityFee: maxPriorityFee,
		gasLimit:       gasLimit,
	}, nil

}

// Submit GGP price
func (t *submitGgpPrice) run() error {

	// Wait for eth client to sync
	if err := services.WaitEthClientSynced(t.c, true); err != nil {
		return err
	}

	// Get node account
	nodeAccount, err := t.w.GetNodeAccount()
	if err != nil {
		return err
	}

	// Data
	var wg errgroup.Group
	var nodeTrusted bool
	var submitPricesEnabled bool

	// Get data
	wg.Go(func() error {
		var err error
		nodeTrusted, err = trustednode.GetMemberExists(t.rp, nodeAccount.Address, nil)
		return err
	})
	wg.Go(func() error {
		var err error
		submitPricesEnabled, err = protocol.GetSubmitPricesEnabled(t.rp, nil)
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return err
	}

	// Check node trusted status & settings
	if !(nodeTrusted && submitPricesEnabled) {
		return nil
	}

	// Log
	t.log.Println("Checking for GGP price checkpoint...")

	// Get block to submit price for
	blockNumber, err := t.getLatestReportableBlock()
	if err != nil {
		return err
	}

	// Allow some blocks to pass in case of a short reorg
	currentBlockNumber, err := t.ec.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	if blockNumber+SubmitFollowDistancePrices > currentBlockNumber {
		return nil
	}

	// Check if a submission needs to be made
	pricesBlock, err := network.GetPricesBlock(t.rp, nil)
	if err != nil {
		return err
	}
	if blockNumber <= pricesBlock {
		return nil
	}

	// If confirm distance has passed, we just want to ensure we have submitted and then early exit
	if blockNumber+ConfirmDistancePrices <= currentBlockNumber {
		hasSubmitted, err := t.hasSubmittedBlockPrices(nodeAccount.Address, blockNumber)
		if err != nil {
			return err
		}
		if hasSubmitted {
			return nil
		}
	}

	// Log
	t.log.Printlnf("Getting GGP price for block %d...", blockNumber)

	// Get GGP price at block
	ggpPrice, err := t.getGgpPrice(blockNumber)
	if err != nil {
		return err
	}

	// Calculate the total effective GGP stake on the network
	zero := new(big.Int).SetUint64(0)
	effectiveGgpStake, err := node.CalculateTotalEffectiveGGPStake(t.rp, zero, zero, ggpPrice, nil)
	if err != nil {
		return fmt.Errorf("Error getting total effective GGP stake: %w", err)
	}

	// Log
	t.log.Printlnf("GGP price: %.6f ETH", mathutils.RoundDown(eth.WeiToEth(ggpPrice), 6))

	// Check if we have reported these specific values before
	hasSubmittedSpecific, err := t.hasSubmittedSpecificBlockPrices(nodeAccount.Address, blockNumber, ggpPrice, effectiveGgpStake)
	if err != nil {
		return err
	}
	if hasSubmittedSpecific {
		return nil
	}

	// We haven't submitted these values, check if we've submitted any for this block so we can log it
	hasSubmitted, err := t.hasSubmittedBlockPrices(nodeAccount.Address, blockNumber)
	if err != nil {
		return err
	}
	if hasSubmitted {
		t.log.Printlnf("Have previously submitted out-of-date prices for block $d, trying again...", blockNumber)
	}

	// Log
	t.log.Println("Submitting GGP price...")

	// Submit GGP price
	if err := t.submitGgpPrice(blockNumber, ggpPrice, effectiveGgpStake); err != nil {
		return fmt.Errorf("Could not submit GGP price: %w", err)
	}

	// Return
	return nil

}

// Get the latest block number to report GGP price for
func (t *submitGgpPrice) getLatestReportableBlock() (uint64, error) {

	// Require eth client synced
	if err := services.RequireEthClientSynced(t.c); err != nil {
		return 0, err
	}

	latestBlock, err := network.GetLatestReportablePricesBlock(t.rp, nil)
	if err != nil {
		return 0, fmt.Errorf("Error getting latest reportable block: %w", err)
	}
	return latestBlock.Uint64(), nil

}

// Check whether prices for a block has already been submitted by the node
func (t *submitGgpPrice) hasSubmittedBlockPrices(nodeAddress common.Address, blockNumber uint64) (bool, error) {

	blockNumberBuf := make([]byte, 32)
	big.NewInt(int64(blockNumber)).FillBytes(blockNumberBuf)
	return t.rp.RocketStorage.GetBool(nil, crypto.Keccak256Hash([]byte("network.prices.submitted.node"), nodeAddress.Bytes(), blockNumberBuf))

}

// Check whether specific prices for a block has already been submitted by the node
func (t *submitGgpPrice) hasSubmittedSpecificBlockPrices(nodeAddress common.Address, blockNumber uint64, ggpPrice, effectiveGgpStake *big.Int) (bool, error) {

	blockNumberBuf := make([]byte, 32)
	big.NewInt(int64(blockNumber)).FillBytes(blockNumberBuf)

	ggpPriceBuf := make([]byte, 32)
	ggpPrice.FillBytes(ggpPriceBuf)

	effectiveGgpStakeBuf := make([]byte, 32)
	effectiveGgpStake.FillBytes(effectiveGgpStakeBuf)

	return t.rp.RocketStorage.GetBool(nil, crypto.Keccak256Hash([]byte("network.prices.submitted.node"), nodeAddress.Bytes(), blockNumberBuf, ggpPriceBuf, effectiveGgpStakeBuf))

}

// Get GGP price at block
func (t *submitGgpPrice) getGgpPrice(blockNumber uint64) (*big.Int, error) {

	// Require 1inch oracle contract
	if err := services.RequireOneInchOracle(t.c); err != nil {
		return nil, err
	}

	// Get GGP token address
	ggpAddress := common.HexToAddress(t.cfg.Rocketpool.GgpTokenAddress)

	// Initialize call options
	opts := &bind.CallOpts{
		BlockNumber: big.NewInt(int64(blockNumber)),
	}

	// Get GGP price
	ggpPrice, err := t.oio.GetRateToEth(opts, ggpAddress, true)
	if err != nil {
		return nil, fmt.Errorf("Could not get GGP price at block %d: %w", blockNumber, err)
	}

	// Return
	return ggpPrice, nil

}

// Submit GGP price and total effective GGP stake
func (t *submitGgpPrice) submitGgpPrice(blockNumber uint64, ggpPrice, effectiveGgpStake *big.Int) error {

	// Log
	t.log.Printlnf("Submitting GGP price for block %d...", blockNumber)

	// Get transactor
	opts, err := t.w.GetNodeAccountTransactor()
	if err != nil {
		return err
	}

	// Get the gas limit
	gasInfo, err := network.EstimateSubmitPricesGas(t.rp, blockNumber, ggpPrice, effectiveGgpStake, opts)
	if err != nil {
		return fmt.Errorf("Could not estimate the gas required to submit GGP price: %w", err)
	}
	var gas *big.Int
	if t.gasLimit != 0 {
		gas = new(big.Int).SetUint64(t.gasLimit)
	} else {
		gas = new(big.Int).SetUint64(gasInfo.SafeGasLimit)
	}

	// Get the max fee
	maxFee := t.maxFee
	if maxFee == nil || maxFee.Uint64() == 0 {
		maxFee, err = rpgas.GetHeadlessMaxFeeWei()
		if err != nil {
			return err
		}
	}

	// Print the gas info
	if !api.PrintAndCheckGasInfo(gasInfo, false, 0, t.log, maxFee, t.gasLimit) {
		return nil
	}

	opts.GasFeeCap = maxFee
	opts.GasTipCap = t.maxPriorityFee
	opts.GasLimit = gas.Uint64()

	// Submit GGP price
	hash, err := network.SubmitPrices(t.rp, blockNumber, ggpPrice, effectiveGgpStake, opts)
	if err != nil {
		return err
	}

	// Print TX info and wait for it to be mined
	err = api.PrintAndWaitForTransaction(t.cfg, hash, t.rp.Client, t.log)
	if err != nil {
		return err
	}

	// Log
	t.log.Printlnf("Successfully submitted GGP price for block %d.", blockNumber)

	// Return
	return nil

}
