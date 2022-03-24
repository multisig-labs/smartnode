package watchtower

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/rocketpool/watchtower/collectors"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/utils/log"
)

// Config
var minTasksInterval, _ = time.ParseDuration("4m")
var maxTasksInterval, _ = time.ParseDuration("6m")
var taskCooldown, _ = time.ParseDuration("10s")

const (
	MaxConcurrentEth1Requests = 200

	RespondChallengesColor           = color.FgWhite
	ClaimGgpRewardsColor             = color.FgGreen
	SubmitGgpPriceColor              = color.FgYellow
	SubmitNetworkBalancesColor       = color.FgYellow
	SubmitWithdrawableMinipoolsColor = color.FgBlue
	DissolveTimedOutMinipoolsColor   = color.FgMagenta
	ProcessWithdrawalsColor          = color.FgCyan
	SubmitScrubMinipoolsColor        = color.FgHiGreen
	ErrorColor                       = color.FgRed
	MetricsColor                     = color.FgHiYellow
)

// Register watchtower command
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Run Rocket Pool watchtower activity daemon",
		Action: func(c *cli.Context) error {
			return run(c)
		},
	})
}

// Run daemon
func run(c *cli.Context) error {

	// Configure
	configureHTTP()

	// Wait until node is registered
	if err := services.WaitNodeRegistered(c, true); err != nil {
		return err
	}

	// Initialize the scrub metrics reporter
	scrubCollector := collectors.NewScrubCollector()

	// Initialize tasks
	respondChallenges, err := newRespondChallenges(c, log.NewColorLogger(RespondChallengesColor))
	if err != nil {
		return err
	}
	claimGgpRewards, err := newClaimGgpRewards(c, log.NewColorLogger(ClaimGgpRewardsColor))
	if err != nil {
		return err
	}
	submitGgpPrice, err := newSubmitGgpPrice(c, log.NewColorLogger(SubmitGgpPriceColor))
	if err != nil {
		return err
	}
	submitNetworkBalances, err := newSubmitNetworkBalances(c, log.NewColorLogger(SubmitNetworkBalancesColor))
	if err != nil {
		return err
	}
	submitWithdrawableMinipools, err := newSubmitWithdrawableMinipools(c, log.NewColorLogger(SubmitWithdrawableMinipoolsColor))
	if err != nil {
		return err
	}
	dissolveTimedOutMinipools, err := newDissolveTimedOutMinipools(c, log.NewColorLogger(DissolveTimedOutMinipoolsColor))
	if err != nil {
		return err
	}
	processWithdrawals, err := newProcessWithdrawals(c, log.NewColorLogger(ProcessWithdrawalsColor))
	if err != nil {
		return err
	}
	submitScrubMinipools, err := newSubmitScrubMinipools(c, log.NewColorLogger(SubmitScrubMinipoolsColor), scrubCollector)
	if err != nil {
		return err
	}

	// Initialize error logger
	errorLog := log.NewColorLogger(ErrorColor)

	intervalDelta := maxTasksInterval - minTasksInterval
	secondsDelta := intervalDelta.Seconds()

	// Wait group to handle the various threads
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// Run task loop
	go func() {
		for {
			// Randomize the next interval
			randomSeconds := rand.Intn(int(secondsDelta))
			interval := time.Duration(randomSeconds)*time.Second + minTasksInterval

			if err := respondChallenges.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(taskCooldown)
			if err := claimGgpRewards.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(taskCooldown)
			if err := submitGgpPrice.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(taskCooldown)
			if err := submitNetworkBalances.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(taskCooldown)
			if err := submitWithdrawableMinipools.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(taskCooldown)
			if err := dissolveTimedOutMinipools.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(taskCooldown)
			if err := processWithdrawals.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(taskCooldown)
			if err := submitScrubMinipools.run(); err != nil {
				errorLog.Println(err)
			}
			time.Sleep(interval)
		}
		wg.Done()
	}()

	// Run metrics loop
	go func() {
		err := runMetricsServer(c, log.NewColorLogger(MetricsColor), scrubCollector)
		if err != nil {
			errorLog.Println(err)
		}
		wg.Done()
	}()

	// Wait for both threads to stop
	wg.Wait()
	return nil
}

// Configure HTTP transport settings
func configureHTTP() {

	// The watchtower daemon makes a large number of concurrent RPC requests to the Eth1 client
	// The HTTP transport is set to cache connections for future re-use equal to the maximum expected number of concurrent requests
	// This prevents issues related to memory consumption and address allowance from repeatedly opening and closing connections
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = MaxConcurrentEth1Requests

}
