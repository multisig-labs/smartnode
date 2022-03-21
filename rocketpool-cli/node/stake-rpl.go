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

func nodeStakeGgp(c *cli.Context) error {

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Get node status
	status, err := rp.NodeStatus()
	if err != nil {
		return err
	}

	// If a custom nonce is set, print the multi-transaction warning
	if c.GlobalUint64("nonce") != 0 {
		cliutils.PrintMultiTransactionNonceWarning()
	}

	// Check for fixed-supply GGP balance
	ggpBalance := *(status.AccountBalances.GGP)
	if status.AccountBalances.FixedSupplyGGP.Cmp(big.NewInt(0)) > 0 {

		// Confirm swapping GGP
		if c.Bool("swap") || cliutils.Confirm(fmt.Sprintf("The node has a balance of %.6f old GGP. Would you like to swap it for new GGP before staking?", math.RoundDown(eth.WeiToEth(status.AccountBalances.FixedSupplyGGP), 6))) {

			// Check allowance
			allowance, err := rp.GetNodeSwapGgpAllowance()
			if err != nil {
				return err
			}

			if allowance.Allowance.Cmp(status.AccountBalances.FixedSupplyGGP) < 0 {
				fmt.Println("Before swapping legacy GGP for new GGP, you must first give the new GGP contract approval to interact with your legacy GGP.")
				fmt.Println("This only needs to be done once for your node.")

				// If a custom nonce is set, print the multi-transaction warning
				if c.GlobalUint64("nonce") != 0 {
					cliutils.PrintMultiTransactionNonceWarning()
				}

				// Calculate max uint256 value
				maxApproval := big.NewInt(2)
				maxApproval = maxApproval.Exp(maxApproval, big.NewInt(256), nil)
				maxApproval = maxApproval.Sub(maxApproval, big.NewInt(1))

				// Get approval gas
				approvalGas, err := rp.NodeSwapGgpApprovalGas(maxApproval)
				if err != nil {
					return err
				}
				// Assign max fees
				err = gas.AssignMaxFeeAndLimit(approvalGas.GasInfo, rp, c.Bool("yes"))
				if err != nil {
					return err
				}

				// Prompt for confirmation
				if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf("Do you want to let the new GGP contract interact with your legacy GGP?"))) {
					fmt.Println("Cancelled.")
					return nil
				}

				// Approve GGP for swapping
				response, err := rp.NodeSwapGgpApprove(maxApproval)
				if err != nil {
					return err
				}
				hash := response.ApproveTxHash
				fmt.Printf("Approving legacy GGP for swapping...\n")
				cliutils.PrintTransactionHash(rp, hash)
				if _, err = rp.WaitForTransaction(hash); err != nil {
					return err
				}
				fmt.Println("Successfully approved access to legacy GGP.")

				// If a custom nonce is set, increment it for the next transaction
				if c.GlobalUint64("nonce") != 0 {
					rp.IncrementCustomNonce()
				}
			}

			// Check GGP can be swapped
			canSwap, err := rp.CanNodeSwapGgp(status.AccountBalances.FixedSupplyGGP)
			if err != nil {
				return err
			}
			if !canSwap.CanSwap {
				fmt.Println("Cannot swap GGP:")
				if canSwap.InsufficientBalance {
					fmt.Println("The node's old GGP balance is insufficient.")
				}
				return nil
			}
			fmt.Println("GGP Swap Gas Info:")
			// Assign max fees
			err = gas.AssignMaxFeeAndLimit(canSwap.GasInfo, rp, c.Bool("yes"))
			if err != nil {
				return err
			}

			// Prompt for confirmation
			if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf("Are you sure you want to swap %.6f old GGP for new GGP?", math.RoundDown(eth.WeiToEth(status.AccountBalances.FixedSupplyGGP), 6)))) {
				fmt.Println("Cancelled.")
				return nil
			}

			// Swap GGP
			swapResponse, err := rp.NodeSwapGgp(status.AccountBalances.FixedSupplyGGP)
			if err != nil {
				return err
			}

			fmt.Printf("Swapping old GGP for new GGP...\n")
			cliutils.PrintTransactionHash(rp, swapResponse.SwapTxHash)
			if _, err = rp.WaitForTransaction(swapResponse.SwapTxHash); err != nil {
				return err
			}

			// Log
			fmt.Printf("Successfully swapped %.6f old GGP for new GGP.\n", math.RoundDown(eth.WeiToEth(status.AccountBalances.FixedSupplyGGP), 6))
			fmt.Println("")

			// If a custom nonce is set, increment it for the next transaction
			if c.GlobalUint64("nonce") != 0 {
				rp.IncrementCustomNonce()
			}

			// Get new account GGP balance
			ggpBalance.Add(status.AccountBalances.GGP, status.AccountBalances.FixedSupplyGGP)

		}

	}

	// Get stake mount
	var amountWei *big.Int
	if c.String("amount") == "min" {

		// Set amount to min per minipool GGP stake
		ggpPrice, err := rp.GgpPrice()
		if err != nil {
			return err
		}
		amountWei = ggpPrice.MinPerMinipoolGgpStake

	} else if c.String("amount") == "max" {

		// Set amount to max per minipool GGP stake
		ggpPrice, err := rp.GgpPrice()
		if err != nil {
			return err
		}
		amountWei = ggpPrice.MaxPerMinipoolGgpStake

	} else if c.String("amount") == "all" {

		// Set amount to node's entire GGP balance
		amountWei = &ggpBalance

	} else if c.String("amount") != "" {

		// Parse amount
		stakeAmount, err := strconv.ParseFloat(c.String("amount"), 64)
		if err != nil {
			return fmt.Errorf("Invalid stake amount '%s': %w", c.String("amount"), err)
		}
		amountWei = eth.EthToWei(stakeAmount)

	} else {

		// Get min/max per minipool GGP stake amounts
		ggpPrice, err := rp.GgpPrice()
		if err != nil {
			return err
		}
		minAmount := ggpPrice.MinPerMinipoolGgpStake
		maxAmount := ggpPrice.MaxPerMinipoolGgpStake

		// Prompt for amount option
		amountOptions := []string{
			fmt.Sprintf("The minimum minipool stake amount (%.6f GGP)?", math.RoundUp(eth.WeiToEth(minAmount), 6)),
			fmt.Sprintf("The maximum effective minipool stake amount (%.6f GGP)?", math.RoundUp(eth.WeiToEth(maxAmount), 6)),
			fmt.Sprintf("Your entire GGP balance (%.6f GGP)?", math.RoundDown(eth.WeiToEth(&ggpBalance), 6)),
			"A custom amount",
		}
		selected, _ := cliutils.Select("Please choose an amount of GGP to stake:", amountOptions)
		switch selected {
		case 0:
			amountWei = minAmount
		case 1:
			amountWei = maxAmount
		case 2:
			amountWei = &ggpBalance
		}

		// Prompt for custom amount
		if amountWei == nil {
			inputAmount := cliutils.Prompt("Please enter an amount of GGP to stake:", "^\\d+(\\.\\d+)?$", "Invalid amount")
			stakeAmount, err := strconv.ParseFloat(inputAmount, 64)
			if err != nil {
				return fmt.Errorf("Invalid stake amount '%s': %w", inputAmount, err)
			}
			amountWei = eth.EthToWei(stakeAmount)
		}

	}

	// Check allowance
	allowance, err := rp.GetNodeStakeGgpAllowance()
	if err != nil {
		return err
	}

	if allowance.Allowance.Cmp(amountWei) < 0 {
		fmt.Println("Before staking GGP, you must first give the staking contract approval to interact with your GGP.")
		fmt.Println("This only needs to be done once for your node.")

		// If a custom nonce is set, print the multi-transaction warning
		if c.GlobalUint64("nonce") != 0 {
			cliutils.PrintMultiTransactionNonceWarning()
		}

		// Calculate max uint256 value
		maxApproval := big.NewInt(2)
		maxApproval = maxApproval.Exp(maxApproval, big.NewInt(256), nil)
		maxApproval = maxApproval.Sub(maxApproval, big.NewInt(1))

		// Get approval gas
		approvalGas, err := rp.NodeStakeGgpApprovalGas(maxApproval)
		if err != nil {
			return err
		}
		// Assign max fees
		err = gas.AssignMaxFeeAndLimit(approvalGas.GasInfo, rp, c.Bool("yes"))
		if err != nil {
			return err
		}

		// Prompt for confirmation
		if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf("Do you want to let the staking contract interact with your GGP?"))) {
			fmt.Println("Cancelled.")
			return nil
		}

		// Approve GGP for staking
		response, err := rp.NodeStakeGgpApprove(maxApproval)
		if err != nil {
			return err
		}
		hash := response.ApproveTxHash
		fmt.Printf("Approving GGP for staking...\n")
		cliutils.PrintTransactionHash(rp, hash)
		if _, err = rp.WaitForTransaction(hash); err != nil {
			return err
		}
		fmt.Println("Successfully approved staking access to GGP.")

		// If a custom nonce is set, increment it for the next transaction
		if c.GlobalUint64("nonce") != 0 {
			rp.IncrementCustomNonce()
		}
	}

	// Check GGP can be staked
	canStake, err := rp.CanNodeStakeGgp(amountWei)
	if err != nil {
		return err
	}
	if !canStake.CanStake {
		fmt.Println("Cannot stake GGP:")
		if canStake.InsufficientBalance {
			fmt.Println("The node's GGP balance is insufficient.")
		}
		if !canStake.InConsensus {
			fmt.Println("The GGP price and total effective staked GGP of the network are still being voted on by the Oracle DAO.\nPlease try again in a few minutes.")
		}
		return nil
	}

	fmt.Println("GGP Stake Gas Info:")
	// Assign max fees
	err = gas.AssignMaxFeeAndLimit(canStake.GasInfo, rp, c.Bool("yes"))
	if err != nil {
		return err
	}

	// Prompt for confirmation
	if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf("Are you sure you want to stake %.6f GGP? You will not be able to unstake this GGP until you exit your validators and close your minipools, or reach over 150%% collateral!", math.RoundDown(eth.WeiToEth(amountWei), 6)))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Stake GGP
	stakeResponse, err := rp.NodeStakeGgp(amountWei)
	if err != nil {
		return err
	}

	fmt.Printf("Staking GGP...\n")
	cliutils.PrintTransactionHash(rp, stakeResponse.StakeTxHash)
	if _, err = rp.WaitForTransaction(stakeResponse.StakeTxHash); err != nil {
		return err
	}

	// Log & return
	fmt.Printf("Successfully staked %.6f GGP.\n", math.RoundDown(eth.WeiToEth(amountWei), 6))
	return nil

}
