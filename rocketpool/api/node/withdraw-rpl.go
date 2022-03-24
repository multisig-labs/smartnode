package node

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
)

func canNodeWithdrawGgp(c *cli.Context, amountWei *big.Int) (*api.CanNodeWithdrawGgpResponse, error) {

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

	// Response
	response := api.CanNodeWithdrawGgpResponse{}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Data
	var wg errgroup.Group
	var ggpStake *big.Int
	var minimumGgpStake *big.Int
	var currentTime uint64
	var ggpStakedTime uint64
	var withdrawalDelay uint64

	// Get GGP stake
	wg.Go(func() error {
		var err error
		ggpStake, err = node.GetNodeGGPStake(rp, nodeAccount.Address, nil)
		return err
	})

	// Get minimum GGP stake
	wg.Go(func() error {
		var err error
		minimumGgpStake, err = node.GetNodeMinimumGGPStake(rp, nodeAccount.Address, nil)
		return err
	})

	// Get current block
	wg.Go(func() error {
		header, err := ec.HeaderByNumber(context.Background(), nil)
		if err == nil {
			currentTime = header.Time
		}
		return err
	})

	// Get GGP staked time
	wg.Go(func() error {
		var err error
		ggpStakedTime, err = node.GetNodeGGPStakedTime(rp, nodeAccount.Address, nil)
		return err
	})

	// Get withdrawal delay
	wg.Go(func() error {
		var err error
		withdrawalDelay, err = protocol.GetRewardsClaimIntervalTime(rp, nil)
		return err
	})

	// Check network consensus
	inConsensus, err := network.InConsensus(rp, nil)
	if err != nil {
		return nil, err
	}
	response.InConsensus = inConsensus

	// Get gas estimate
	wg.Go(func() error {
		opts, err := w.GetNodeAccountTransactor()
		if err != nil {
			return err
		}
		gasInfo, err := node.EstimateWithdrawGGPGas(rp, amountWei, opts)
		if err == nil {
			response.GasInfo = gasInfo
		}
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Check data
	var remainingGgpStake big.Int
	remainingGgpStake.Sub(ggpStake, amountWei)
	response.InsufficientBalance = (amountWei.Cmp(ggpStake) > 0)
	response.MinipoolsUndercollateralized = (remainingGgpStake.Cmp(minimumGgpStake) < 0)
	response.WithdrawalDelayActive = ((currentTime - ggpStakedTime) < withdrawalDelay)

	// Update & return response
	response.CanWithdraw = !(response.InsufficientBalance || response.MinipoolsUndercollateralized || response.WithdrawalDelayActive || !response.InConsensus)
	return &response, nil

}

func nodeWithdrawGgp(c *cli.Context, amountWei *big.Int) (*api.NodeWithdrawGgpResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NodeWithdrawGgpResponse{}

	// Get transactor
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}

	// Override the provided pending TX if requested
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}

	// Withdraw GGP
	hash, err := node.WithdrawGGP(rp, amountWei, opts)
	if err != nil {
		return nil, err
	}
	response.TxHash = hash

	// Return response
	return &response, nil

}
