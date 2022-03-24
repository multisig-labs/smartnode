package collectors

import (
	"context"
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"golang.org/x/sync/errgroup"
)

// Represents the collector for the Performance metrics
type PerformanceCollector struct {
	// The ETH utilization rate (%)
	ethUtilizationRate *prometheus.Desc

	// The total amount of ETH staked
	totalStakingBalanceEth *prometheus.Desc

	// The ETH / ggpAVAX ratio
	ethGgpavaxExchangeRate *prometheus.Desc

	// The total amount of ETH locked (TVL)
	totalValueLockedEth *prometheus.Desc

	// The total ggpAVAX supply
	totalGgpavaxSupply *prometheus.Desc

	// The ETH balance of the ggpAVAX contract address
	ggpavaxContractBalance *prometheus.Desc

	// The Rocket Pool contract manager
	rp *rocketpool.RocketPool
}

// Create a new PerformanceCollector instance
func NewPerformanceCollector(rp *rocketpool.RocketPool) *PerformanceCollector {
	subsystem := "performance"
	return &PerformanceCollector{
		ethUtilizationRate: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "eth_utilization_rate"),
			"The ETH utilization rate (%)",
			nil, nil,
		),
		totalStakingBalanceEth: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "total_staking_balance_eth"),
			"The total amount of ETH staked",
			nil, nil,
		),
		ethGgpavaxExchangeRate: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "eth_ggpavax_exchange_rate"),
			"The ETH / ggpAVAX ratio",
			nil, nil,
		),
		totalValueLockedEth: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "total_value_locked_eth"),
			"The total amount of ETH locked (TVL)",
			nil, nil,
		),
		ggpavaxContractBalance: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "ggpavax_contract_balance"),
			"The ETH balance of the ggpAVAX contract address",
			nil, nil,
		),
		totalGgpavaxSupply: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "total_ggpavax_supply"),
			"The total ggpAVAX supply",
			nil, nil,
		),
		rp: rp,
	}
}

// Write metric descriptions to the Prometheus channel
func (collector *PerformanceCollector) Describe(channel chan<- *prometheus.Desc) {
	channel <- collector.ethUtilizationRate
	channel <- collector.totalStakingBalanceEth
	channel <- collector.ethGgpavaxExchangeRate
	channel <- collector.totalValueLockedEth
	channel <- collector.ggpavaxContractBalance
	channel <- collector.totalGgpavaxSupply
}

// Collect the latest metric values and pass them to Prometheus
func (collector *PerformanceCollector) Collect(channel chan<- prometheus.Metric) {

	// Sync
	var wg errgroup.Group
	ethUtilizationRate := float64(-1)
	balanceFloat := float64(-1)
	exchangeRate := float64(-1)
	tvlFloat := float64(-1)
	ggpAVAXBalance := float64(-1)
	ggpavaxFloat := float64(-1)

	// Get the ETH utilization rate
	wg.Go(func() error {
		_ethUtilizationRate, err := network.GetETHUtilizationRate(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting ETH utilization rate: %w", err)
		} else {
			ethUtilizationRate = _ethUtilizationRate
		}
		return nil
	})

	// Get the total ETH staking balance
	wg.Go(func() error {
		totalStakingBalance, err := network.GetStakingETHBalance(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting total ETH staking balance: %w", err)
		} else {
			balanceFloat = eth.WeiToEth(totalStakingBalance)
		}
		return nil
	})

	// Get the ETH-ggpAVAX exchange rate
	wg.Go(func() error {
		_exchangeRate, err := tokens.GetGGPAVAXExchangeRate(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting ETH-ggpAVAX exchange rate: %w", err)
		} else {
			exchangeRate = _exchangeRate
		}
		return nil
	})

	// Get the total ETH balance (TVL)
	wg.Go(func() error {
		tvl, err := network.GetTotalETHBalance(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting total ETH balance (TVL): %w", err)
		} else {
			tvlFloat = eth.WeiToEth(tvl)
		}
		return nil
	})

	// Get the ETH balance of the ggpAVAX contract
	wg.Go(func() error {
		ggpAVAXContract, err := collector.rp.GetContract("gogoTokenGGPAVAX")
		if err != nil {
			return fmt.Errorf("Error getting ETH balance of ggpAVAX staking contract: %w", err)
		} else {
			balance, err := collector.rp.Client.BalanceAt(context.Background(), *ggpAVAXContract.Address, nil)
			if err != nil {
				return fmt.Errorf("Error getting ETH balance of ggpAVAX staking contract: %w", err)
			} else {
				ggpAVAXBalance = eth.WeiToEth(balance)
			}
		}
		return nil
	})

	// Get the total ggpAVAX supply
	wg.Go(func() error {
		totalGgpavaxSupply, err := tokens.GetGGPAVAXTotalSupply(collector.rp, nil)
		if err != nil {
			return fmt.Errorf("Error getting total ggpAVAX supply: %w", err)
		} else {
			ggpavaxFloat = eth.WeiToEth(totalGgpavaxSupply)
		}
		return nil
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		log.Printf("%s\n", err.Error())
		return
	}

	channel <- prometheus.MustNewConstMetric(
		collector.ethUtilizationRate, prometheus.GaugeValue, ethUtilizationRate)
	channel <- prometheus.MustNewConstMetric(
		collector.totalStakingBalanceEth, prometheus.GaugeValue, balanceFloat)
	channel <- prometheus.MustNewConstMetric(
		collector.ethGgpavaxExchangeRate, prometheus.GaugeValue, exchangeRate)
	channel <- prometheus.MustNewConstMetric(
		collector.totalValueLockedEth, prometheus.GaugeValue, tvlFloat)
	channel <- prometheus.MustNewConstMetric(
		collector.ggpavaxContractBalance, prometheus.GaugeValue, ggpAVAXBalance)
	channel <- prometheus.MustNewConstMetric(
		collector.totalGgpavaxSupply, prometheus.GaugeValue, ggpavaxFloat)

}
