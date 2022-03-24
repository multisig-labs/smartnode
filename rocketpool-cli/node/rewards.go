package node

import (
	"fmt"
	"time"

	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

func getRewards(c *cli.Context) error {

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Get node GGP rewards status
	rewards, err := rp.NodeRewards()
	if err != nil {
		return err
	}

	if !rewards.Registered {
		fmt.Printf("This node is not currently registered.\n")
		return nil
	}

	colorReset := "\033[0m"
	colorYellow := "\033[33m"

	fmt.Println("=== ETH ===")
	fmt.Printf("You have earned %.4f ETH from the Beacon Chain (including your commissions) so far.\n", rewards.BeaconRewards)

	nextRewardsTime := rewards.LastCheckpoint.Add(rewards.RewardsInterval)
	nextRewardsTimeString := cliutils.GetDateTimeString(uint64(nextRewardsTime.Unix()))
	timeToCheckpointString := time.Until(nextRewardsTime).Round(time.Second).String()
	docsUrl := "https://docs.rocketpool.net/guides/node/rewards.html#claiming-ggp-rewards"

	// Assume 365 days in a year, 24 hours per day
	ggpApr := rewards.EstimatedRewards / rewards.TotalGgpStake / rewards.RewardsInterval.Hours() * (24 * 365) * 100

	fmt.Println("\n=== GGP ===")
	fmt.Printf("The current rewards cycle started on %s.\n", cliutils.GetDateTimeString(uint64(rewards.LastCheckpoint.Unix())))
	fmt.Printf("It will end on %s (%s from now).\n", nextRewardsTimeString, timeToCheckpointString)

	if rewards.UnclaimedRewards > 0 {
		fmt.Printf("%s**WARNING**: you currently have %f GGP unclaimed from the previous cycle. If you don't claim them before the above date, you will lose them!%s\n",
			colorYellow, rewards.UnclaimedRewards, colorReset)
	}
	if rewards.UnclaimedTrustedRewards > 0 {
		fmt.Printf("%s**WARNING**: you currently have %f GGP unclaimed from the previous cycle's Oracle DAO duties. If you don't claim them before the above date, you will lose them!%s\n",
			colorYellow, rewards.UnclaimedTrustedRewards, colorReset)
	}

	fmt.Println()
	fmt.Printf("Your estimated GGP staking rewards for this cycle: %f GGP (this may change based on network activity).\n", rewards.EstimatedRewards)
	fmt.Printf("Based on your current total stake of %f GGP, this is approximately %.2f%% APR.\n", rewards.TotalGgpStake, ggpApr)
	fmt.Printf("Your node has received %f GGP staking rewards in total.\n", rewards.CumulativeRewards)

	if rewards.Trusted {
		ggpTrustedApr := rewards.EstimatedTrustedRewards / rewards.TrustedGgpBond / rewards.RewardsInterval.Hours() * (24 * 365) * 100

		fmt.Println()
		fmt.Printf("You will receive an estimated %f GGP in rewards for Oracle DAO duties (this may change based on network activity).\n", rewards.EstimatedTrustedRewards)
		fmt.Printf("Based on your bond of %f GGP, this is approximately %.2f%% APR.\n", rewards.TrustedGgpBond, ggpTrustedApr)
		fmt.Printf("Your node has received %f GGP Oracle DAO rewards in total.\n", rewards.CumulativeTrustedRewards)
	}

	fmt.Println()
	fmt.Println("These rewards will be claimed automatically when the checkpoint ends, unless you have disabled auto-claims.")
	fmt.Printf("Refer to the Claiming Node Operator Rewards guide at %s for more information.", docsUrl)

	// Return
	return nil

}
