package node

import (
	"bytes"

	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
)

func getStatus(c *cli.Context) (*api.NodeStatusResponse, error) {

	// Get services
	if err := services.RequireNodeWallet(c); err != nil {
		return nil, err
	}
	if err := services.RequireRocketStorage(c); err != nil {
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
	response := api.NodeStatusResponse{}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}
	response.AccountAddress = nodeAccount.Address

	// Sync
	var wg errgroup.Group

	// Get node trusted status
	wg.Go(func() error {
		trusted, err := trustednode.GetMemberExists(rp, nodeAccount.Address, nil)
		if err == nil {
			response.Trusted = trusted
		}
		return err
	})

	// Get node details
	wg.Go(func() error {
		details, err := node.GetNodeDetails(rp, nodeAccount.Address, nil)
		if err == nil {
			response.Registered = details.Exists
			response.WithdrawalAddress = details.WithdrawalAddress
			response.PendingWithdrawalAddress = details.PendingWithdrawalAddress
			response.TimezoneLocation = details.TimezoneLocation
		}
		return err
	})

	// Get node account balances
	wg.Go(func() error {
		var err error
		response.AccountBalances, err = tokens.GetBalances(rp, nodeAccount.Address, nil)
		return err
	})

	// Get staking details
	wg.Go(func() error {
		var err error
		response.GgpStake, err = node.GetNodeGGPStake(rp, nodeAccount.Address, nil)
		return err
	})
	wg.Go(func() error {
		var err error
		response.EffectiveGgpStake, err = node.GetNodeEffectiveGGPStake(rp, nodeAccount.Address, nil)
		return err
	})
	wg.Go(func() error {
		var err error
		response.MinimumGgpStake, err = node.GetNodeMinimumGGPStake(rp, nodeAccount.Address, nil)
		return err
	})
	wg.Go(func() error {
		var err error
		response.MaximumGgpStake, err = node.GetNodeMaximumGGPStake(rp, nodeAccount.Address, nil)
		return err
	})
	wg.Go(func() error {
		var err error
		response.MinipoolLimit, err = node.GetNodeMinipoolLimit(rp, nodeAccount.Address, nil)
		return err
	})

	// Get node minipool counts
	wg.Go(func() error {
		details, err := getNodeMinipoolCountDetails(rp, nodeAccount.Address)
		if err == nil {
			response.MinipoolCounts.Total = len(details)
			for _, mpDetails := range details {
				if mpDetails.Finalised {
					response.MinipoolCounts.Finalised++
				} else {
					switch mpDetails.Status {
					case types.Initialized:
						response.MinipoolCounts.Initialized++
					case types.Prelaunch:
						response.MinipoolCounts.Prelaunch++
					case types.Staking:
						response.MinipoolCounts.Staking++
					case types.Withdrawable:
						response.MinipoolCounts.Withdrawable++
					case types.Dissolved:
						response.MinipoolCounts.Dissolved++
					}
					if mpDetails.RefundAvailable {
						response.MinipoolCounts.RefundAvailable++
					}
					if mpDetails.WithdrawalAvailable {
						response.MinipoolCounts.WithdrawalAvailable++
					}
					if mpDetails.CloseAvailable {
						response.MinipoolCounts.CloseAvailable++
					}
				}
			}
		}
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Get withdrawal address balances
	if !bytes.Equal(nodeAccount.Address.Bytes(), response.WithdrawalAddress.Bytes()) {
		withdrawalBalances, err := tokens.GetBalances(rp, response.WithdrawalAddress, nil)
		if err != nil {
			return nil, err
		}
		response.WithdrawalBalances = withdrawalBalances
	}

	// Get the collateral ratio
	ggpPrice, err := network.GetGGPPrice(rp, nil)
	if err != nil {
		return nil, err
	}
	activeMinipools := response.MinipoolCounts.Total - response.MinipoolCounts.Finalised
	if activeMinipools > 0 {
		response.CollateralRatio = eth.WeiToEth(ggpPrice) * eth.WeiToEth(response.GgpStake) / (float64(activeMinipools) * 16.0)
	} else {
		response.CollateralRatio = -1
	}

	// Return response
	return &response, nil

}
