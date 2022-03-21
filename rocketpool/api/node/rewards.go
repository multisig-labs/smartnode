package node

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/beacon"
	"github.com/rocket-pool/smartnode/shared/types/api"
	apiutils "github.com/rocket-pool/smartnode/shared/utils/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth2"
)

func getRewards(c *cli.Context) (*api.NodeRewardsResponse, error) {

	// Get services
	if err := services.RequireNodeWallet(c); err != nil {
		return nil, err
	}
	if err := services.RequireRocketStorage(c); err != nil {
		return nil, err
	}
	if err := services.RequireEthClientSynced(c); err != nil {
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
	bc, err := services.GetBeaconClient(c)
	if err != nil {
		return nil, err
	}
	cfg, err := services.GetConfig(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NodeRewardsResponse{}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Get the event log interval
	eventLogInterval, err := apiutils.GetEventLogInterval(cfg)
	if err != nil {
		return nil, err
	}

	var totalEffectiveStake *big.Int
	var totalGgpSupply *big.Int
	var inflationInterval *big.Int
	var odaoSize uint64
	var nodeOperatorRewardsPercent float64
	var trustedNodeOperatorRewardsPercent float64
	var totalDepositBalance float64
	var totalNodeShare float64
	var addresses []common.Address
	var beaconHead beacon.BeaconHead

	// Sync
	var wg errgroup.Group

	// Check if the node is registered or not
	wg.Go(func() error {
		exists, err := node.GetNodeExists(rp, nodeAccount.Address, nil)
		if err == nil {
			response.Registered = exists
		}
		return err
	})

	// Get the node registration time
	wg.Go(func() error {
		time, err := rewards.GetNodeRegistrationTime(rp, nodeAccount.Address, nil)
		if err == nil {
			response.NodeRegistrationTime = time
		}
		return err
	})

	// Get node trusted status
	wg.Go(func() error {
		trusted, err := trustednode.GetMemberExists(rp, nodeAccount.Address, nil)
		if err == nil {
			response.Trusted = trusted
		}
		return err
	})

	// Get cumulative rewards
	wg.Go(func() error {
		rewards, err := rewards.CalculateLifetimeNodeRewards(rp, nodeAccount.Address, eventLogInterval, nil)
		if err == nil {
			response.CumulativeRewards = eth.WeiToEth(rewards)
		}
		return err
	})

	// Get the start of the rewards checkpoint
	wg.Go(func() error {
		lastCheckpoint, err := rewards.GetClaimIntervalTimeStart(rp, nil)
		if err == nil {
			response.LastCheckpoint = lastCheckpoint
		}
		return err
	})

	// Get the rewards checkpoint interval
	wg.Go(func() error {
		rewardsInterval, err := rewards.GetClaimIntervalTime(rp, nil)
		if err == nil {
			response.RewardsInterval = rewardsInterval
		}
		return err
	})

	// Get the node's effective stake
	wg.Go(func() error {
		effectiveStake, err := node.GetNodeEffectiveGGPStake(rp, nodeAccount.Address, nil)
		if err == nil {
			response.EffectiveGgpStake = eth.WeiToEth(effectiveStake)
		}
		return err
	})

	// Get the node's total stake
	wg.Go(func() error {
		stake, err := node.GetNodeGGPStake(rp, nodeAccount.Address, nil)
		if err == nil {
			response.TotalGgpStake = eth.WeiToEth(stake)
		}
		return err
	})

	// Get the total network effective stake
	wg.Go(func() error {
		totalEffectiveStake, err = node.GetTotalEffectiveGGPStake(rp, nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Get the total GGP supply
	wg.Go(func() error {
		totalGgpSupply, err = tokens.GetGGPTotalSupply(rp, nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Get the GGP inflation interval
	wg.Go(func() error {
		inflationInterval, err = tokens.GetGGPInflationIntervalRate(rp, nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Get the node operator rewards percent
	wg.Go(func() error {
		nodeOperatorRewardsPercent, err = rewards.GetNodeOperatorRewardsPercent(rp, nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Check if rewards are currently available from the previous checkpoint
	wg.Go(func() error {
		unclaimedRewardsWei, err := rewards.GetNodeClaimRewardsAmount(rp, nodeAccount.Address, nil)
		if err == nil {
			response.UnclaimedRewards = eth.WeiToEth(unclaimedRewardsWei)
		}
		return err
	})

	// Get the list of minipool addresses for this node
	wg.Go(func() error {
		_addresses, err := minipool.GetNodeMinipoolAddresses(rp, nodeAccount.Address, nil)
		if err != nil {
			return fmt.Errorf("Error getting node minipool addresses: %w", err)
		}
		addresses = _addresses
		return nil
	})

	// Get the beacon head
	wg.Go(func() error {
		_beaconHead, err := bc.GetBeaconHead()
		if err != nil {
			return fmt.Errorf("Error getting beacon chain head: %w", err)
		}
		beaconHead = _beaconHead
		return nil
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Calculate the total deposits and corresponding beacon chain balance share
	minipoolDetails, err := eth2.GetBeaconBalances(rp, bc, addresses, beaconHead, nil)
	if err != nil {
		return nil, err
	}
	for _, minipool := range minipoolDetails {
		totalDepositBalance += eth.WeiToEth(minipool.NodeDeposit)
		totalNodeShare += eth.WeiToEth(minipool.NodeBalance)
	}
	response.BeaconRewards = totalNodeShare - totalDepositBalance

	// Calculate the estimated rewards
	rewardsIntervalDays := response.RewardsInterval.Seconds() / (60 * 60 * 24)
	inflationPerDay := eth.WeiToEth(inflationInterval)
	totalGgpAtNextCheckpoint := (math.Pow(inflationPerDay, float64(rewardsIntervalDays)) - 1) * eth.WeiToEth(totalGgpSupply)
	if totalGgpAtNextCheckpoint < 0 {
		totalGgpAtNextCheckpoint = 0
	}

	if totalEffectiveStake.Cmp(big.NewInt(0)) == 1 {
		response.EstimatedRewards = response.EffectiveGgpStake / eth.WeiToEth(totalEffectiveStake) * totalGgpAtNextCheckpoint * nodeOperatorRewardsPercent
	}

	if response.Trusted {

		var wg2 errgroup.Group

		// Get the node registration time
		wg2.Go(func() error {
			time, err := rewards.GetTrustedNodeRegistrationTime(rp, nodeAccount.Address, nil)
			if err == nil {
				response.TrustedNodeRegistrationTime = time
			}
			return err
		})

		// Get cumulative ODAO rewards
		wg2.Go(func() error {
			rewards, err := rewards.CalculateLifetimeTrustedNodeRewards(rp, nodeAccount.Address, eventLogInterval, nil)
			if err == nil {
				response.CumulativeTrustedRewards = eth.WeiToEth(rewards)
			}
			return err
		})

		// Get the ODAO member count
		wg2.Go(func() error {
			odaoSize, err = trustednode.GetMemberCount(rp, nil)
			if err != nil {
				return err
			}
			return nil
		})

		// Get the trusted node operator rewards percent
		wg2.Go(func() error {
			trustedNodeOperatorRewardsPercent, err = rewards.GetTrustedNodeOperatorRewardsPercent(rp, nil)
			if err != nil {
				return err
			}
			return nil
		})

		// Get the node's oDAO GGP stake
		wg2.Go(func() error {
			bond, err := trustednode.GetMemberGGPBondAmount(rp, nodeAccount.Address, nil)
			if err == nil {
				response.TrustedGgpBond = eth.WeiToEth(bond)
			}
			return err
		})

		// Check if rewards are currently available from the previous checkpoint for the ODAO
		wg2.Go(func() error {
			unclaimedRewardsWei, err := rewards.GetTrustedNodeClaimRewardsAmount(rp, nodeAccount.Address, nil)
			if err == nil {
				response.UnclaimedTrustedRewards = eth.WeiToEth(unclaimedRewardsWei)
			}
			return err
		})

		// Wait for data
		if err := wg2.Wait(); err != nil {
			return nil, err
		}

		response.EstimatedTrustedRewards = totalGgpAtNextCheckpoint * trustedNodeOperatorRewardsPercent / float64(odaoSize)

	}

	// Return response
	return &response, nil

}
