package network

import (
	"github.com/rocket-pool/rocketpool-go/deposit"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
)

func getStats(c *cli.Context) (*api.NetworkStatsResponse, error) {

	// Get services
	if err := services.RequireRocketStorage(c); err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NetworkStatsResponse{}

	// Sync
	var wg errgroup.Group

	// Get the deposit pool balance
	wg.Go(func() error {
		balance, err := deposit.GetBalance(rp, nil)
		if err == nil {
			response.DepositPoolBalance = eth.WeiToEth(balance)
		}
		return err
	})

	// Get the total minipool capacity
	wg.Go(func() error {
		minipoolQueueCapacity, err := minipool.GetQueueCapacity(rp, nil)
		if err == nil {
			response.MinipoolCapacity = eth.WeiToEth(minipoolQueueCapacity.Total)
		}
		return err
	})

	// Get the ETH utilization rate
	wg.Go(func() error {
		stakerUtilization, err := network.GetETHUtilizationRate(rp, nil)
		if err == nil {
			response.StakerUtilization = stakerUtilization
		}
		return err
	})

	// Get node fee
	wg.Go(func() error {
		nodeFee, err := network.GetNodeFee(rp, nil)
		if err == nil {
			response.NodeFee = nodeFee
		}
		return err
	})

	// Get node count
	wg.Go(func() error {
		nodeCount, err := node.GetNodeCount(rp, nil)
		if err == nil {
			response.NodeCount = nodeCount
		}
		return err
	})

	// Get minipool counts
	wg.Go(func() error {
		minipoolCounts, err := minipool.GetMinipoolCountPerStatus(rp, nil)
		if err != nil {
			return err
		}
		response.InitializedMinipoolCount = minipoolCounts.Initialized.Uint64()
		response.PrelaunchMinipoolCount = minipoolCounts.Prelaunch.Uint64()
		response.StakingMinipoolCount = minipoolCounts.Staking.Uint64()
		response.WithdrawableMinipoolCount = minipoolCounts.Withdrawable.Uint64()
		response.DissolvedMinipoolCount = minipoolCounts.Dissolved.Uint64()

		finalizedCount, err := minipool.GetFinalisedMinipoolCount(rp, nil)
		if err != nil {
			return err
		}
		response.FinalizedMinipoolCount = finalizedCount

		return nil
	})

	// Get GGP price
	wg.Go(func() error {
		ggpPrice, err := network.GetGGPPrice(rp, nil)
		if err == nil {
			response.GgpPrice = eth.WeiToEth(ggpPrice)
		}
		return err
	})

	// Get total GGP staked
	wg.Go(func() error {
		totalStaked, err := node.GetTotalGGPStake(rp, nil)
		if err == nil {
			response.TotalGgpStaked = eth.WeiToEth(totalStaked)
		}
		return err
	})

	// Get total effective GGP staked
	wg.Go(func() error {
		effectiveStaked, err := node.GetTotalEffectiveGGPStake(rp, nil)
		if err == nil {
			response.EffectiveGgpStaked = eth.WeiToEth(effectiveStaked)
		}
		return err
	})

	// Get ggpAVAX price
	wg.Go(func() error {
		ggpavaxPrice, err := tokens.GetGGPAVAXExchangeRate(rp, nil)
		if err == nil {
			response.GgpavaxPrice = ggpavaxPrice
		}
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Get the TVL
	activeMinipools := response.InitializedMinipoolCount +
		response.PrelaunchMinipoolCount +
		response.StakingMinipoolCount +
		response.WithdrawableMinipoolCount +
		response.DissolvedMinipoolCount
	tvl := float64(activeMinipools)*32 + response.DepositPoolBalance + response.MinipoolCapacity + (response.TotalGgpStaked * response.GgpPrice)
	response.TotalValueLocked = tvl

	// Return response
	return &response, nil

}
