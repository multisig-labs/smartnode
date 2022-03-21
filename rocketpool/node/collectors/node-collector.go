package collectors

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/smartnode/shared/services/beacon"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/utils/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth2"
	"golang.org/x/sync/errgroup"
)

// Represents the collector for the user's node
type NodeCollector struct {
	// The total amount of GGP staked on the node
	totalStakedGgp *prometheus.Desc

	// The effective amount of GGP staked on the node (honoring the 150% collateral cap)
	effectiveStakedGgp *prometheus.Desc

	// The GGP collateral level for the node
	ggpCollateral *prometheus.Desc

	// The cumulative GGP rewards earned by the node
	cumulativeGgpRewards *prometheus.Desc

	// The expected GGP rewards for the node at the next rewards checkpoint
	expectedGgpRewards *prometheus.Desc

	// The estimated APR of GGP for the node from the next rewards checkpoint
	ggpApr *prometheus.Desc

	// The token balances of your node wallet
	balances *prometheus.Desc

	// The number of active minipools owned by the node
	activeMinipoolCount *prometheus.Desc

	// The amount of ETH this node deposited into minipools
	depositedEth *prometheus.Desc

	// The node's total share of its minipool's beacon chain balances
	beaconShare *prometheus.Desc

	// The total balances of all this node's validators on the beacon chain
	beaconBalance *prometheus.Desc

	// The GGP rewards from the last period that have not been claimed yet
	unclaimedRewards *prometheus.Desc

	// The Rocket Pool contract manager
	rp *rocketpool.RocketPool

	// The beacon client
	bc beacon.Client

	// The node's address
	nodeAddress common.Address

	// The event log interval for the current eth1 client
	eventLogInterval *big.Int

	// The next block to start from when looking at cumulative GGP rewards
	nextRewardsStartBlock *big.Int

	// The cumulative amount of GGP earned
	cumulativeRewards float64
}

// Create a new NodeCollector instance
func NewNodeCollector(rp *rocketpool.RocketPool, bc beacon.Client, nodeAddress common.Address, cfg config.RocketPoolConfig) *NodeCollector {

	// Get the event log interval
	eventLogInterval, err := api.GetEventLogInterval(cfg)
	if err != nil {
		log.Printf("Error getting event log interval: %s\n", err.Error())
		return nil
	}

	subsystem := "node"
	return &NodeCollector{
		totalStakedGgp: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "total_staked_ggp"),
			"The total amount of GGP staked on the node",
			nil, nil,
		),
		effectiveStakedGgp: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "effective_staked_ggp"),
			"The effective amount of GGP staked on the node (honoring the 150% collateral cap)",
			nil, nil,
		),
		ggpCollateral: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "ggp_collateral"),
			"The GGP collateral level for the node",
			nil, nil,
		),
		cumulativeGgpRewards: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "cumulative_ggp_rewards"),
			"The cumulative GGP rewards earned by the node",
			nil, nil,
		),
		expectedGgpRewards: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "expected_ggp_rewards"),
			"The expected GGP rewards for the node at the next rewards checkpoint",
			nil, nil,
		),
		ggpApr: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "ggp_apr"),
			"The estimated APR of GGP for the node from the next rewards checkpoint",
			nil, nil,
		),
		balances: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "balance"),
			"How much ETH is in this node wallet",
			[]string{"Token"}, nil,
		),
		activeMinipoolCount: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "active_minipool_count"),
			"The number of active minipools owned by the node",
			nil, nil,
		),
		depositedEth: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "deposited_eth"),
			"The amount of ETH this node deposited into minipools",
			nil, nil,
		),
		beaconShare: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "beacon_share"),
			"The node's total share of its minipool's beacon chain balances",
			nil, nil,
		),
		beaconBalance: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "beacon_balance"),
			"The total balances of all this node's validators on the beacon chain",
			nil, nil,
		),
		unclaimedRewards: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "unclaimed_rewards"),
			"The GGP rewards from the last period that have not been claimed yet",
			nil, nil,
		),
		rp:               rp,
		bc:               bc,
		nodeAddress:      nodeAddress,
		eventLogInterval: eventLogInterval,
	}
}

// Write metric descriptions to the Prometheus channel
func (collector *NodeCollector) Describe(channel chan<- *prometheus.Desc) {
	channel <- collector.totalStakedGgp
	channel <- collector.effectiveStakedGgp
	channel <- collector.cumulativeGgpRewards
	channel <- collector.expectedGgpRewards
	channel <- collector.ggpApr
	channel <- collector.balances
	channel <- collector.activeMinipoolCount
	channel <- collector.depositedEth
	channel <- collector.beaconShare
	channel <- collector.unclaimedRewards
}

// Collect the latest metric values and pass them to Prometheus
func (collector *NodeCollector) Collect(channel chan<- prometheus.Metric) {

	// Sync
	var wg errgroup.Group
	stakedGgp := float64(-1)
	effectiveStakedGgp := float64(-1)
	var rewardsInterval time.Duration
	var inflationInterval *big.Int
	var totalGgpSupply *big.Int
	var totalEffectiveStake *big.Int
	var nodeOperatorRewardsPercent float64
	ethBalance := float64(-1)
	oldGgpBalance := float64(-1)
	newGgpBalance := float64(-1)
	ggpavaxBalance := float64(-1)
	var activeMinipoolCount float64
	var ggpPrice float64
	collateralRatio := float64(-1)
	var addresses []common.Address
	var beaconHead beacon.BeaconHead
	unclaimedRewards := float64(-1)

	// Get the total staked GGP
	wg.Go(func() error {
		stakedGgpWei, err := node.GetNodeGGPStake(collector.rp, collector.nodeAddress, nil)
		if err != nil {
			return fmt.Errorf("Error getting total staked GGP: %w", err)
		} else {
			stakedGgp = eth.WeiToEth(stakedGgpWei)
		}
		return nil
	})

	// Get the effective staked GGP
	wg.Go(func() error {
		effectiveStakedGgpWei, err := node.GetNodeEffectiveGGPStake(collector.rp, collector.nodeAddress, nil)
		if err != nil {
			return fmt.Errorf("Error getting effective staked GGP: %w", err)
		} else {
			effectiveStakedGgp = eth.WeiToEth(effectiveStakedGgpWei)
		}
		return nil
	})

	// Get the cumulative GGP rewards
	wg.Go(func() error {
		cumulativeRewardsWei, err := rewards.CalculateLifetimeNodeRewards(collector.rp, collector.nodeAddress, collector.eventLogInterval, collector.nextRewardsStartBlock)
		if err != nil {
			return fmt.Errorf("Error getting cumulative GGP rewards: %w", err)
		}

		header, err := collector.rp.Client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("Error getting latest block header: %w", err)
		}

		collector.cumulativeRewards += eth.WeiToEth(cumulativeRewardsWei)
		collector.nextRewardsStartBlock = big.NewInt(0).Add(header.Number, big.NewInt(1))
		return nil
	})

	// Get the rewards checkpoint interval
	wg.Go(func() error {
		_rewardsInterval, err := rewards.GetClaimIntervalTime(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting rewards checkpoint interval: %w", err)
		}
		rewardsInterval = _rewardsInterval
		return nil
	})

	// Get the GGP inflation interval
	wg.Go(func() error {
		_inflationInterval, err := tokens.GetGGPInflationIntervalRate(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting GGP inflation interval: %w", err)
		}
		inflationInterval = _inflationInterval
		return nil
	})

	// Get the total GGP supply
	wg.Go(func() error {
		_totalGgpSupply, err := tokens.GetGGPTotalSupply(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting total GGP supply: %w", err)
		}
		totalGgpSupply = _totalGgpSupply
		return nil
	})

	// Get the total network effective stake
	wg.Go(func() error {
		_totalEffectiveStake, err := node.GetTotalEffectiveGGPStake(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting total network effective stake: %w", err)
		}
		totalEffectiveStake = _totalEffectiveStake
		return nil
	})

	// Get the node operator rewards percent
	wg.Go(func() error {
		_nodeOperatorRewardsPercent, err := rewards.GetNodeOperatorRewardsPercent(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting node operator rewards percent: %w", err)
		}
		nodeOperatorRewardsPercent = _nodeOperatorRewardsPercent
		return nil
	})

	// Get the node balances
	wg.Go(func() error {
		balances, err := tokens.GetBalances(collector.rp, collector.nodeAddress, nil)
		if err != nil {
			return fmt.Errorf("Error getting node balances: %w", err)
		}
		ethBalance = eth.WeiToEth(balances.ETH)
		oldGgpBalance = eth.WeiToEth(balances.FixedSupplyGGP)
		newGgpBalance = eth.WeiToEth(balances.GGP)
		ggpavaxBalance = eth.WeiToEth(balances.GGPAVAX)
		return nil
	})

	// Get the number of active minipools on the node
	wg.Go(func() error {
		_activeMinipoolCount, err := minipool.GetNodeActiveMinipoolCount(collector.rp, collector.nodeAddress, nil)
		if err != nil {
			return fmt.Errorf("Error getting node active minipool count: %w", err)
		}
		activeMinipoolCount = float64(_activeMinipoolCount)
		return nil
	})

	// Get the GGP price
	wg.Go(func() error {
		ggpPriceWei, err := network.GetGGPPrice(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting GGP price: %w", err)
		}
		ggpPrice = eth.WeiToEth(ggpPriceWei)
		return nil
	})

	// Get the list of minipool addresses for this node
	wg.Go(func() error {
		_addresses, err := minipool.GetNodeMinipoolAddresses(collector.rp, collector.nodeAddress, nil)
		if err != nil {
			return fmt.Errorf("Error getting node minipool addresses: %w", err)
		}
		addresses = _addresses
		return nil
	})

	// Get the beacon head
	wg.Go(func() error {
		_beaconHead, err := collector.bc.GetBeaconHead()
		if err != nil {
			return fmt.Errorf("Error getting beacon chain head: %w", err)
		}
		beaconHead = _beaconHead
		return nil
	})

	// Get the GGP price
	wg.Go(func() error {
		unclaimedRewardsWei, err := rewards.GetNodeClaimRewardsAmount(collector.rp, collector.nodeAddress, nil)
		if err != nil {
			return fmt.Errorf("Error getting GGP price: %w", err)
		}
		unclaimedRewards = eth.WeiToEth(unclaimedRewardsWei)
		return nil
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	// Calculate the estimated rewards
	rewardsIntervalDays := rewardsInterval.Seconds() / (60 * 60 * 24)
	inflationPerDay := eth.WeiToEth(inflationInterval)
	totalGgpAtNextCheckpoint := (math.Pow(inflationPerDay, float64(rewardsIntervalDays)) - 1) * eth.WeiToEth(totalGgpSupply)
	if totalGgpAtNextCheckpoint < 0 {
		totalGgpAtNextCheckpoint = 0
	}
	estimatedRewards := float64(0)
	if totalEffectiveStake.Cmp(big.NewInt(0)) == 1 {
		estimatedRewards = effectiveStakedGgp / eth.WeiToEth(totalEffectiveStake) * totalGgpAtNextCheckpoint * nodeOperatorRewardsPercent
	}

	// Calculate the GGP APR
	ggpApr := estimatedRewards / stakedGgp / rewardsInterval.Hours() * (24 * 365) * 100

	// Calculate the collateral ratio
	if activeMinipoolCount > 0 {
		collateralRatio = ggpPrice * stakedGgp / (activeMinipoolCount * 16.0)
	}

	// Calculate the total deposits and corresponding beacon chain balance share
	minipoolDetails, err := eth2.GetBeaconBalances(collector.rp, collector.bc, addresses, beaconHead, nil)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}
	totalDepositBalance := float64(0)
	totalNodeShare := float64(0)
	totalBeaconBalance := float64(0)
	for _, minipool := range minipoolDetails {
		totalDepositBalance += eth.WeiToEth(minipool.NodeDeposit)
		totalNodeShare += eth.WeiToEth(minipool.NodeBalance)
		totalBeaconBalance += eth.WeiToEth(minipool.TotalBalance)
	}

	// Update all the metrics
	channel <- prometheus.MustNewConstMetric(
		collector.totalStakedGgp, prometheus.GaugeValue, stakedGgp)
	channel <- prometheus.MustNewConstMetric(
		collector.effectiveStakedGgp, prometheus.GaugeValue, effectiveStakedGgp)
	channel <- prometheus.MustNewConstMetric(
		collector.ggpCollateral, prometheus.GaugeValue, collateralRatio)
	channel <- prometheus.MustNewConstMetric(
		collector.cumulativeGgpRewards, prometheus.GaugeValue, collector.cumulativeRewards)
	channel <- prometheus.MustNewConstMetric(
		collector.expectedGgpRewards, prometheus.GaugeValue, estimatedRewards)
	channel <- prometheus.MustNewConstMetric(
		collector.ggpApr, prometheus.GaugeValue, ggpApr)
	channel <- prometheus.MustNewConstMetric(
		collector.balances, prometheus.GaugeValue, ethBalance, "ETH")
	channel <- prometheus.MustNewConstMetric(
		collector.balances, prometheus.GaugeValue, oldGgpBalance, "Legacy GGP")
	channel <- prometheus.MustNewConstMetric(
		collector.balances, prometheus.GaugeValue, newGgpBalance, "New GGP")
	channel <- prometheus.MustNewConstMetric(
		collector.balances, prometheus.GaugeValue, ggpavaxBalance, "ggpAVAX")
	channel <- prometheus.MustNewConstMetric(
		collector.activeMinipoolCount, prometheus.GaugeValue, activeMinipoolCount)
	channel <- prometheus.MustNewConstMetric(
		collector.depositedEth, prometheus.GaugeValue, totalDepositBalance)
	channel <- prometheus.MustNewConstMetric(
		collector.beaconShare, prometheus.GaugeValue, totalNodeShare)
	channel <- prometheus.MustNewConstMetric(
		collector.beaconBalance, prometheus.GaugeValue, totalBeaconBalance)
	channel <- prometheus.MustNewConstMetric(
		collector.unclaimedRewards, prometheus.GaugeValue, unclaimedRewards)
}
