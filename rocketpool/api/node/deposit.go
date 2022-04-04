package node

import (
	"context"
	"crypto/x509"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	tndao "github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"github.com/rocket-pool/rocketpool-go/settings/trustednode"
	tnsettings "github.com/rocket-pool/rocketpool-go/settings/trustednode"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/beacon"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
	"github.com/rocket-pool/smartnode/shared/utils/validator"
)

type minipoolCreated struct {
	Minipool common.Address
	Node     common.Address
	Time     *big.Int
}

func canNodeDeposit(c *cli.Context, amountWei *big.Int, minNodeFee float64, salt *big.Int) (*api.CanNodeDepositResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
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
	bc, err := services.GetBeaconClient(c)
	if err != nil {
		return nil, err
	}

	// Get eth2 config
	eth2Config, err := bc.GetEth2Config()
	if err != nil {
		return nil, err
	}

	// Response
	response := api.CanNodeDepositResponse{}

	// Check if amount is zero
	amountIsZero := (amountWei.Cmp(big.NewInt(0)) == 0)

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Adjust the salt
	if salt.Cmp(big.NewInt(0)) == 0 {
		nonce, err := ec.NonceAt(context.Background(), nodeAccount.Address, nil)
		if err != nil {
			return nil, err
		}
		salt.SetUint64(nonce)
	}

	// Data
	var wg1 errgroup.Group
	var isTrusted bool
	var minipoolCount uint64
	var minipoolLimit uint64
	var minipoolAddress common.Address

	// Check node balance
	wg1.Go(func() error {
		ethBalanceWei, err := ec.BalanceAt(context.Background(), nodeAccount.Address, nil)
		if err == nil {
			response.InsufficientBalance = (amountWei.Cmp(ethBalanceWei) > 0)
		}
		return err
	})

	// Check node deposits are enabled
	wg1.Go(func() error {
		depositEnabled, err := protocol.GetNodeDepositEnabled(rp, nil)
		if err == nil {
			response.DepositDisabled = !depositEnabled
		}
		return err
	})

	// Get trusted status
	wg1.Go(func() error {
		var err error
		isTrusted, err = tndao.GetMemberExists(rp, nodeAccount.Address, nil)
		return err
	})

	// Get node staking information
	wg1.Go(func() error {
		var err error
		minipoolCount, err = minipool.GetNodeMinipoolCount(rp, nodeAccount.Address, nil)
		return err
	})
	wg1.Go(func() error {
		var err error
		minipoolLimit, err = node.GetNodeMinipoolLimit(rp, nodeAccount.Address, nil)
		return err
	})

	// Get consensus status
	wg1.Go(func() error {
		var err error
		inConsensus, err := network.InConsensus(rp, nil)
		response.InConsensus = inConsensus
		return err
	})

	// Get gas estimate
	wg1.Go(func() error {
		opts, err := w.GetNodeAccountTransactor()
		if err != nil {
			return err
		}
		opts.Value = amountWei

		// Get the deposit type
		depositType, err := node.GetDepositType(rp, amountWei, nil)
		if err != nil {
			return err
		}

		// Get the next validator key
		validatorKey, err := validator.GetValidatorPrivateKey("/home/chandler/.gogopool/data/validator/staking/staking.key") // this is just temporary
		if err != nil {
			return err
		}

		// Get the next minipool address and withdrawal credentials
		minipoolAddress, err = utils.GenerateAddress(rp, nodeAccount.Address, depositType, salt, nil)
		if err != nil {
			return err
		}
		withdrawalCredentials, err := minipool.GetMinipoolWithdrawalCredentials(rp, minipoolAddress, nil)
		if err != nil {
			return err
		}

		// Get validator deposit data and associated parameters
		depositData, depositDataRoot, err := validator.GetDepositData(validatorKey, withdrawalCredentials, eth2Config)
		if err != nil {
			return err
		}
		pubKey := rptypes.BytesToValidatorPubkey(depositData.PublicKey)
		signature := rptypes.BytesToValidatorSignature(depositData.Signature)

		// Run the deposit gas estimator
		gasInfo, err := node.EstimateDepositGas(rp, minNodeFee, pubKey, signature, depositDataRoot, salt, minipoolAddress, opts)
		if err == nil {
			response.GasInfo = gasInfo
		}
		return err
	})

	// Wait for data
	if err := wg1.Wait(); err != nil {
		return nil, err
	}

	// Check data
	response.InsufficientGgpStake = (minipoolCount >= minipoolLimit)
	response.MinipoolAddress = minipoolAddress
	response.InvalidAmount = (!isTrusted && amountIsZero)

	// Check oracle node unbonded minipool limit
	if isTrusted && amountIsZero {

		// Data
		var wg2 errgroup.Group
		var unbondedMinipoolCount uint64
		var unbondedMinipoolsMax uint64

		// Get unbonded minipool details
		wg2.Go(func() error {
			var err error
			unbondedMinipoolCount, err = tndao.GetMemberUnbondedValidatorCount(rp, nodeAccount.Address, nil)
			return err
		})
		wg2.Go(func() error {
			var err error
			unbondedMinipoolsMax, err = tnsettings.GetMinipoolUnbondedMax(rp, nil)
			return err
		})

		// Wait for data
		if err := wg2.Wait(); err != nil {
			return nil, err
		}

		// Check unbonded minipool limit
		response.UnbondedMinipoolsAtMax = (unbondedMinipoolCount >= unbondedMinipoolsMax)

	}

	// Update & return response
	response.CanDeposit = !(response.InsufficientBalance || response.InsufficientGgpStake || response.InvalidAmount || response.UnbondedMinipoolsAtMax || response.DepositDisabled || !response.InConsensus)
	return &response, nil

}

func nodeDeposit(c *cli.Context, amountWei *big.Int, minNodeFee float64, salt *big.Int) (*api.NodeDepositResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
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
	bc, err := services.GetBeaconClient(c)
	if err != nil {
		return nil, err
	}

	// Get eth2 config
	eth2Config, err := bc.GetEth2Config()
	if err != nil {
		return nil, err
	}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NodeDepositResponse{}

	// Adjust the salt
	if salt.Cmp(big.NewInt(0)) == 0 {
		nonce, err := ec.NonceAt(context.Background(), nodeAccount.Address, nil)
		if err != nil {
			return nil, err
		}
		salt.SetUint64(nonce)
	}

	// Get the scrub period
	scrubPeriodUnix, err := trustednode.GetScrubPeriod(rp, nil)
	if err != nil {
		return nil, err
	}
	scrubPeriod := time.Duration(scrubPeriodUnix) * time.Second
	response.ScrubPeriod = scrubPeriod

	// Get transactor
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	opts.Value = amountWei

	// Get the deposit type
	depositType, err := node.GetDepositType(rp, amountWei, nil)
	if err != nil {
		return nil, err
	}

	// Create and save a new validator key
	// TODO: move staking key to config
	validatorKey, err := validator.GetValidatorPrivateKey("/Users/pisti/GolandProjects/avalanchego/staking/local/staker3.key") // this is just temporary
	if err != nil {
		return nil, err
	}
	fmt.Println(validatorKey.PublicKey)
	// Get the next minipool address and withdrawal credentials
	minipoolAddress, err := utils.GenerateAddress(rp, nodeAccount.Address, depositType, salt, nil)
	if err != nil {
		return nil, err
	}

	// convert nodeAccount.Address to bytes
	nodeAccountAddressBytes := nodeAccount.Address.Bytes()
	// convert it to a common Hash
	withdrawalCredentials := common.BytesToHash(nodeAccountAddressBytes)

	// Get validator deposit data and associated parameters
	depositData, _, err := validator.GetDepositData(validatorKey, withdrawalCredentials, eth2Config)
	if err != nil {
		return nil, err
	}
	pubKey := rptypes.BytesToValidatorPubkey(x509.MarshalPKCS1PublicKey(&validatorKey.PublicKey))
	_ = rptypes.BytesToValidatorSignature(depositData.Signature)

	// Make sure a validator with this pubkey doesn't already exist
	status := beacon.ValidatorStatus{
		Pubkey:                     rptypes.BytesToValidatorPubkey(depositData.PublicKey),
		Index:                      0,
		WithdrawalCredentials:      withdrawalCredentials,
		Balance:                    0,
		EffectiveBalance:           0,
		Slashed:                    false,
		ActivationEligibilityEpoch: 0,
		ActivationEpoch:            0,
		ExitEpoch:                  0,
		WithdrawableEpoch:          0,
		Exists:                     false,
	}
	//status, err := bc.GetValidatorStatus(pubKey, nil)
	if err != nil {
		return nil, fmt.Errorf("Error checking for existing validator status: %w\nYour funds have not been deposited for your own safety.", err)
	}
	if status.Exists {
		return nil, fmt.Errorf("**** ALERT ****\n"+
			"Your minipool %s has the following as a validator pubkey:\n\t%s\n"+
			"This key is already in use by validator %d on the Beacon chain!\n"+
			"Rocket Pool will not allow you to deposit this validator for your own safety so you do not get slashed.\n"+
			"PLEASE REPORT THIS TO THE ROCKET POOL DEVELOPERS.\n"+
			"***************\n", minipoolAddress.Hex(), pubKey.Hex(), status.Index)
	}

	// Override the provided pending TX if requested
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}

	// Get avalanche-go NodeId
	nodeIdResponse, err := bc.GetNodeId()
	if err != nil {
		return nil, err
	}
	fmt.Println(nodeIdResponse.NodeID)
	// Deposit
	hash, err := node.Deposit(rp, minNodeFee, pubKey, salt, minipoolAddress, nodeIdResponse.NodeID, opts)
	if err != nil {
		return nil, err
	}

	// Save wallet
	if err := w.Save(); err != nil {
		return nil, err
	}

	response.TxHash = hash
	response.MinipoolAddress = minipoolAddress
	response.ValidatorPubkey = pubKey

	// Return response
	return &response, nil

}

