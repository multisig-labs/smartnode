package odao

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/gas"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)

func join(c *cli.Context) error {

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
	if status.AccountBalances.FixedSupplyGGP.Cmp(big.NewInt(0)) > 0 {

		// Confirm swapping GGP
		if c.Bool("swap") || cliutils.Confirm(fmt.Sprintf("The node has a balance of %.6f old GGP. Would you like to swap it for new GGP before transferring your bond?", math.RoundDown(eth.WeiToEth(status.AccountBalances.FixedSupplyGGP), 6))) {

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
		}
	}

	// Check if node can join the oracle DAO
	canJoin, err := rp.CanJoinTNDAO()
	if err != nil {
		return err
	}
	if !canJoin.CanJoin {
		fmt.Println("Cannot join the oracle DAO:")
		if canJoin.ProposalExpired {
			fmt.Println("The proposal for you to join the oracle DAO does not exist or has expired.")
		}
		if canJoin.AlreadyMember {
			fmt.Println("The node is already a member of the oracle DAO.")
		}
		if canJoin.InsufficientGgpBalance {
			fmt.Println("The node does not have enough GGP to pay the GGP bond.")
		}
		return nil
	}

	// Display gas estimate
	// Assign max fees
	err = gas.AssignMaxFeeAndLimit(canJoin.GasInfo, rp, c.Bool("yes"))
	if err != nil {
		return err
	}
	rp.PrintMultiTxWarning()

	// Prompt for confirmation
	if !(c.Bool("yes") || cliutils.Confirm("Are you sure you want to join the oracle DAO? Your GGP bond will be locked until you leave.")) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Approve GGP for joining the ODAO
	response, err := rp.ApproveGGPToJoinTNDAO()
	if err != nil {
		return err
	}
	hash := response.ApproveTxHash
	fmt.Printf("Approving GGP for joining the Oracle DAO...\n")
	cliutils.PrintTransactionHashNoCancel(rp, hash)

	// If a custom nonce is set, increment it for the next transaction
	if c.GlobalUint64("nonce") != 0 {
		rp.IncrementCustomNonce()
	}

	// Join the ODAO
	joinResponse, err := rp.JoinTNDAO(hash)
	if err != nil {
		return err
	}
	fmt.Printf("Joining the ODAO...\n")
	cliutils.PrintTransactionHash(rp, joinResponse.JoinTxHash)
	if _, err = rp.WaitForTransaction(joinResponse.JoinTxHash); err != nil {
		return err
	}

	// Log & return
	fmt.Println("Successfully joined the oracle DAO.")
	return nil

}
