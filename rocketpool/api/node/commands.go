package node

import (
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/utils/api"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

// Register subcommands
func RegisterSubcommands(command *cli.Command, name string, aliases []string) {
	command.Subcommands = append(command.Subcommands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the node",
		Subcommands: []cli.Command{

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get the node's status",
				UsageText: "rocketpool api node status",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(getStatus(c))
					return nil

				},
			},

			{
				Name:      "sync",
				Aliases:   []string{"y"},
				Usage:     "Get the sync progress of the eth1 and eth2 clients",
				UsageText: "rocketpool api node sync",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(getSyncProgress(c))
					return nil

				},
			},

			{
				Name:      "can-register",
				Usage:     "Check whether the node can be registered with Rocket Pool",
				UsageText: "rocketpool api node can-register timezone-location",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					timezoneLocation, err := cliutils.ValidateTimezoneLocation("timezone location", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canRegisterNode(c, timezoneLocation))
					return nil

				},
			},
			{
				Name:      "register",
				Aliases:   []string{"r"},
				Usage:     "Register the node with Rocket Pool",
				UsageText: "rocketpool api node register timezone-location",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					timezoneLocation, err := cliutils.ValidateTimezoneLocation("timezone location", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(registerNode(c, timezoneLocation))
					return nil

				},
			},

			{
				Name:      "can-set-withdrawal-address",
				Usage:     "Checks if the node can set its withdrawal address",
				UsageText: "rocketpool api node can-set-withdrawal-address address confirm",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					withdrawalAddress, err := cliutils.ValidateAddress("withdrawal address", c.Args().Get(0))
					if err != nil {
						return err
					}

					confirm, err := cliutils.ValidateBool("confirm", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canSetWithdrawalAddress(c, withdrawalAddress, confirm))
					return nil

				},
			},
			{
				Name:      "set-withdrawal-address",
				Aliases:   []string{"w"},
				Usage:     "Set the node's withdrawal address",
				UsageText: "rocketpool api node set-withdrawal-address address confirm",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					withdrawalAddress, err := cliutils.ValidateAddress("withdrawal address", c.Args().Get(0))
					if err != nil {
						return err
					}

					confirm, err := cliutils.ValidateBool("confirm", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(setWithdrawalAddress(c, withdrawalAddress, confirm))
					return nil

				},
			},

			{
				Name:      "can-confirm-withdrawal-address",
				Usage:     "Checks if the node can confirm its withdrawal address",
				UsageText: "rocketpool api node can-confirm-withdrawal-address",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(canConfirmWithdrawalAddress(c))
					return nil

				},
			},
			{
				Name:      "confirm-withdrawal-address",
				Usage:     "Confirms the node's withdrawal address if it was set back to the node address",
				UsageText: "rocketpool api node confirm-withdrawal-address",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(confirmWithdrawalAddress(c))
					return nil

				},
			},

			{
				Name:      "can-set-timezone",
				Usage:     "Checks if the node can set its timezone location",
				UsageText: "rocketpool api node can-set-timezone timezone-location",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					timezoneLocation, err := cliutils.ValidateTimezoneLocation("timezone location", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canSetTimezoneLocation(c, timezoneLocation))
					return nil

				},
			},
			{
				Name:      "set-timezone",
				Aliases:   []string{"t"},
				Usage:     "Set the node's timezone location",
				UsageText: "rocketpool api node set-timezone timezone-location",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					timezoneLocation, err := cliutils.ValidateTimezoneLocation("timezone location", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(setTimezoneLocation(c, timezoneLocation))
					return nil

				},
			},

			{
				Name:      "can-swap-ggp",
				Usage:     "Check whether the node can swap old GGP for new GGP",
				UsageText: "rocketpool api node can-swap-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("swap amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canNodeSwapGgp(c, amountWei))
					return nil

				},
			},
			{
				Name:      "swap-ggp-approve-ggp",
				Aliases:   []string{"p1"},
				Usage:     "Approve fixed-supply GGP for swapping to new GGP",
				UsageText: "rocketpool api node swap-ggp-approve-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("swap amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(approveFsGgp(c, amountWei))
					return nil

				},
			},
			{
				Name:      "wait-and-swap-ggp",
				Aliases:   []string{"p2"},
				Usage:     "Swap old GGP for new GGP, waiting for the approval TX hash to be mined first",
				UsageText: "rocketpool api node wait-and-swap-ggp amount tx-hash",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("swap amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					hash, err := cliutils.ValidateTxHash("swap amount", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(waitForApprovalAndSwapFsGgp(c, amountWei, hash))
					return nil

				},
			},
			{
				Name:      "get-swap-ggp-approval-gas",
				Usage:     "Estimate the gas cost of legacy GGP interaction approval",
				UsageText: "rocketpool api node get-swap-ggp-approval-gas",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("approve amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(getSwapApprovalGas(c, amountWei))
					return nil

				},
			},
			{
				Name:      "swap-ggp-allowance",
				Usage:     "Get the node's legacy GGP allowance for new GGP contract",
				UsageText: "rocketpool api node swap-allowance-ggp",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(allowanceFsGgp(c))
					return nil

				},
			},
			{
				Name:      "swap-ggp",
				Aliases:   []string{"p3"},
				Usage:     "Swap old GGP for new GGP",
				UsageText: "rocketpool api node swap-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("swap amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(swapGgp(c, amountWei))
					return nil

				},
			},

			{
				Name:      "can-stake-ggp",
				Usage:     "Check whether the node can stake GGP",
				UsageText: "rocketpool api node can-stake-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("stake amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canNodeStakeGgp(c, amountWei))
					return nil

				},
			},
			{
				Name:      "stake-ggp-approve-ggp",
				Aliases:   []string{"k1"},
				Usage:     "Approve GGP for staking against the node",
				UsageText: "rocketpool api node stake-ggp-approve-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("stake amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(approveGgp(c, amountWei))
					return nil

				},
			},
			{
				Name:      "wait-and-stake-ggp",
				Aliases:   []string{"k2"},
				Usage:     "Stake GGP against the node, waiting for approval tx-hash to be mined first",
				UsageText: "rocketpool api node wait-and-stake-ggp amount tx-hash",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("stake amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					hash, err := cliutils.ValidateTxHash("tx-hash", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(waitForApprovalAndStakeGgp(c, amountWei, hash))
					return nil

				},
			},
			{
				Name:      "get-stake-ggp-approval-gas",
				Usage:     "Estimate the gas cost of new GGP interaction approval",
				UsageText: "rocketpool api node get-stake-ggp-approval-gas",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("approve amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(getStakeApprovalGas(c, amountWei))
					return nil

				},
			},
			{
				Name:      "stake-ggp-allowance",
				Usage:     "Get the node's GGP allowance for the staking contract",
				UsageText: "rocketpool api node stake-allowance-ggp",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(allowanceGgp(c))
					return nil

				},
			},
			{
				Name:      "stake-ggp",
				Aliases:   []string{"k3"},
				Usage:     "Stake GGP against the node",
				UsageText: "rocketpool api node stake-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("stake amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(stakeGgp(c, amountWei))
					return nil

				},
			},

			{
				Name:      "can-withdraw-ggp",
				Usage:     "Check whether the node can withdraw staked GGP",
				UsageText: "rocketpool api node can-withdraw-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("withdrawal amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canNodeWithdrawGgp(c, amountWei))
					return nil

				},
			},
			{
				Name:      "withdraw-ggp",
				Aliases:   []string{"i"},
				Usage:     "Withdraw GGP staked against the node",
				UsageText: "rocketpool api node withdraw-ggp amount",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("withdrawal amount", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(nodeWithdrawGgp(c, amountWei))
					return nil

				},
			},

			{
				Name:      "can-deposit",
				Usage:     "Check whether the node can make a deposit",
				UsageText: "rocketpool api node can-deposit amount min-fee salt",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 3); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidateDepositWeiAmount("deposit amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					minNodeFee, err := cliutils.ValidateFraction("minimum node fee", c.Args().Get(1))
					if err != nil {
						return err
					}
					salt, err := cliutils.ValidateBigInt("salt", c.Args().Get(2))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canNodeDeposit(c, amountWei, minNodeFee, salt))
					return nil

				},
			},
			{
				Name:      "deposit",
				Aliases:   []string{"d"},
				Usage:     "Make a deposit and create a minipool",
				UsageText: "rocketpool api node deposit amount min-fee salt",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 3); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidateDepositWeiAmount("deposit amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					minNodeFee, err := cliutils.ValidateFraction("minimum node fee", c.Args().Get(1))
					if err != nil {
						return err
					}
					salt, err := cliutils.ValidateBigInt("salt", c.Args().Get(2))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(nodeDeposit(c, amountWei, minNodeFee, salt))
					return nil

				},
			},

			{
				Name:      "can-send",
				Usage:     "Check whether the node can send ETH or tokens to an address",
				UsageText: "rocketpool api node can-send amount token",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("send amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					token, err := cliutils.ValidateTokenType("token type", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canNodeSend(c, amountWei, token))
					return nil

				},
			},
			{
				Name:      "send",
				Aliases:   []string{"n"},
				Usage:     "Send ETH or tokens from the node account to an address",
				UsageText: "rocketpool api node send amount token to",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 3); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("send amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					token, err := cliutils.ValidateTokenType("token type", c.Args().Get(1))
					if err != nil {
						return err
					}
					toAddress, err := cliutils.ValidateAddress("to address", c.Args().Get(2))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(nodeSend(c, amountWei, token, toAddress))
					return nil

				},
			},

			{
				Name:      "can-burn",
				Usage:     "Check whether the node can burn tokens for ETH",
				UsageText: "rocketpool api node can-burn amount token",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("burn amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					token, err := cliutils.ValidateBurnableTokenType("token type", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(canNodeBurn(c, amountWei, token))
					return nil

				},
			},
			{
				Name:      "burn",
				Aliases:   []string{"b"},
				Usage:     "Burn tokens for ETH",
				UsageText: "rocketpool api node burn amount token",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 2); err != nil {
						return err
					}
					amountWei, err := cliutils.ValidatePositiveWeiAmount("burn amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					token, err := cliutils.ValidateBurnableTokenType("token type", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					api.PrintResponse(nodeBurn(c, amountWei, token))
					return nil

				},
			},

			{
				Name:      "can-claim-ggp-rewards",
				Usage:     "Check whether the node has GGP rewards available to claim",
				UsageText: "rocketpool api node can-claim-ggp-rewards",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(canNodeClaimGgp(c))
					return nil

				},
			},
			{
				Name:      "claim-ggp-rewards",
				Usage:     "Claim available GGP rewards",
				UsageText: "rocketpool api node claim-ggp-rewards",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(nodeClaimGgp(c))
					return nil

				},
			},

			{
				Name:      "rewards",
				Usage:     "Get GGP rewards info",
				UsageText: "rocketpool api node rewards",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(getRewards(c))
					return nil

				},
			},

			{
				Name:      "deposit-contract-info",
				Usage:     "Get information about the deposit contract specified by Rocket Pool and the Beacon Chain client",
				UsageText: "rocketpool api node deposit-contract-info",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(getDepositContractInfo(c))
					return nil

				},
			},
		},
	})
}
