package node

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/gas"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)

func nodeWithdrawGgp(c *cli.Context) error {

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Get withdrawal mount
	var amountWei *big.Int
	if c.String("amount") == "max" {

		// Get node status
		status, err := rp.NodeStatus()
		if err != nil {
			return err
		}

		// Set amount to maximum withdrawable amount
		var maxAmount big.Int
		if status.GgpStake.Cmp(status.MinimumGgpStake) > 0 {
			maxAmount.Sub(status.GgpStake, status.MinimumGgpStake)
		}
		amountWei = &maxAmount

	} else if c.String("amount") != "" {

		// Parse amount
		withdrawalAmount, err := strconv.ParseFloat(c.String("amount"), 64)
		if err != nil {
			return fmt.Errorf("Invalid withdrawal amount '%s': %w", c.String("amount"), err)
		}
		amountWei = eth.EthToWei(withdrawalAmount)

	} else {

		// Get node status
		status, err := rp.NodeStatus()
		if err != nil {
			return err
		}

		// Get maximum withdrawable amount
		var maxAmount big.Int
		maxAmount.Sub(status.GgpStake, status.MaximumGgpStake)
		if maxAmount.Sign() == 1 {
			// Prompt for maximum amount
			if cliutils.Confirm(fmt.Sprintf("Would you like to withdraw the maximum amount of staked GGP (%.6f GGP)?", math.RoundDown(eth.WeiToEth(&maxAmount), 6))) {
				amountWei = &maxAmount
			} else {

				// Prompt for custom amount
				inputAmount := cliutils.Prompt("Please enter an amount of staked GGP to withdraw:", "^\\d+(\\.\\d+)?$", "Invalid amount")
				withdrawalAmount, err := strconv.ParseFloat(inputAmount, 64)
				if err != nil {
					return fmt.Errorf("Invalid withdrawal amount '%s': %w", inputAmount, err)
				}
				amountWei = eth.EthToWei(withdrawalAmount)

			}
		} else {
			fmt.Printf("Cannot withdraw staked GGP - you have %.6f GGP staked, but are not allowed to withdraw below %.6f GGP (150%% collateral).\n",
				math.RoundDown(eth.WeiToEth(status.GgpStake), 6),
				math.RoundDown(eth.WeiToEth(status.MaximumGgpStake), 6))
			return nil
		}

	}

	// Check GGP can be withdrawn
	canWithdraw, err := rp.CanNodeWithdrawGgp(amountWei)
	if err != nil {
		return err
	}
	if !canWithdraw.CanWithdraw {
		fmt.Println("Cannot withdraw staked GGP:")
		if canWithdraw.InsufficientBalance {
			fmt.Println("The node's staked GGP balance is insufficient.")
		}
		if canWithdraw.MinipoolsUndercollateralized {
			fmt.Println("Remaining staked GGP is not enough to collateralize the node's minipools.")
		}
		if canWithdraw.WithdrawalDelayActive {
			fmt.Println("The withdrawal delay period has not passed.")
		}
		if !canWithdraw.InConsensus {
			fmt.Println("The GGP price and total effective staked GGP of the network are still being voted on by the Oracle DAO.\nPlease try again in a few minutes.")
		}
		return nil
	}

	// Assign max fees
	err = gas.AssignMaxFeeAndLimit(canWithdraw.GasInfo, rp, c.Bool("yes"))
	if err != nil {
		return err
	}

	// Prompt for confirmation
	if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf("Are you sure you want to withdraw %.6f staked GGP? This may decrease your node's GGP rewards.", math.RoundDown(eth.WeiToEth(amountWei), 6)))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Withdraw GGP
	response, err := rp.NodeWithdrawGgp(amountWei)
	if err != nil {
		return err
	}

	fmt.Printf("Withdrawing GGP...\n")
	cliutils.PrintTransactionHash(rp, response.TxHash)
	if _, err = rp.WaitForTransaction(response.TxHash); err != nil {
		return err
	}

	// Log & return
	fmt.Printf("Successfully withdrew %.6f staked GGP.\n", math.RoundDown(eth.WeiToEth(amountWei), 6))
	return nil

}
