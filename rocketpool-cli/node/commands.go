package node

import (
	"github.com/urfave/cli"

	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the node",
		Subcommands: []cli.Command{

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get the node's status",
				UsageText: "rocketpool node status",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getStatus(c)

				},
			},

			{
				Name:      "sync",
				Aliases:   []string{"y"},
				Usage:     "Get the sync progress of the eth1 and eth2 clients",
				UsageText: "rocketpool node sync",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getSyncProgress(c)

				},
			},

			{
				Name:      "register",
				Aliases:   []string{"r"},
				Usage:     "Register the node with Rocket Pool",
				UsageText: "rocketpool node register [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "timezone, t",
						Usage: "The timezone location to register the node with (in the format 'Country/City')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("timezone") != "" {
						if _, err := cliutils.ValidateTimezoneLocation("timezone location", c.String("timezone")); err != nil {
							return err
						}
					}

					// Run
					return registerNode(c)

				},
			},

			{
				Name:      "rewards",
				Aliases:   []string{"e"},
				Usage:     "Get the time and your expected GGP rewards of the next checkpoint",
				UsageText: "rocketpool node rewards",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getRewards(c)

				},
			},

			{
				Name:      "set-withdrawal-address",
				Aliases:   []string{"w"},
				Usage:     "Set the node's withdrawal address",
				UsageText: "rocketpool node set-withdrawal-address [options] address",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm setting withdrawal address",
					},
					cli.BoolFlag{
						Name:  "force",
						Usage: "Force update the withdrawal address, bypassing the 'pending' state that requires a confirmation transaction from the new address",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					withdrawalAddress, err := cliutils.ValidateAddress("withdrawal address", c.Args().Get(0))
					if err != nil {
						return err
					}

					// Run
					return setWithdrawalAddress(c, withdrawalAddress)

				},
			},

			{
				Name:      "confirm-withdrawal-address",
				Aliases:   []string{"f"},
				Usage:     "Confirm the node's pending withdrawal address if it has been set back to the node's address itself",
				UsageText: "rocketpool node confirm-withdrawal-address [options]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm withdrawal address",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return confirmWithdrawalAddress(c)

				},
			},

			{
				Name:      "set-timezone",
				Aliases:   []string{"t"},
				Usage:     "Set the node's timezone location",
				UsageText: "rocketpool node set-timezone [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "timezone, t",
						Usage: "The timezone location to set for the node (in the format 'Country/City')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("timezone") != "" {
						if _, err := cliutils.ValidateTimezoneLocation("timezone location", c.String("timezone")); err != nil {
							return err
						}
					}

					// Run
					return setTimezoneLocation(c)

				},
			},

			{
				Name:      "swap-ggp",
				Aliases:   []string{"p"},
				Usage:     "Swap old GGP for new GGP",
				UsageText: "rocketpool node swap-ggp [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of old GGP to swap (or 'all')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" && c.String("amount") != "all" {
						if _, err := cliutils.ValidatePositiveEthAmount("swap amount", c.String("amount")); err != nil {
							return err
						}
					}

					// Run
					return nodeSwapGgp(c)

				},
			},

			{
				Name:      "stake-ggp",
				Aliases:   []string{"k"},
				Usage:     "Stake GGP against the node",
				UsageText: "rocketpool node stake-ggp [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of GGP to stake (or 'min', 'max', or 'all')",
					},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm GGP stake",
					},
					cli.BoolFlag{
						Name:  "swap, s",
						Usage: "Automatically confirm swapping old GGP before staking",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" && c.String("amount") != "min" && c.String("amount") != "max" && c.String("amount") != "all" {
						if _, err := cliutils.ValidatePositiveEthAmount("stake amount", c.String("amount")); err != nil {
							return err
						}
					}

					// Run
					return nodeStakeGgp(c)

				},
			},

			{
				Name:      "claim-ggp",
				Aliases:   []string{"c"},
				Usage:     "Claim available GGP rewards for the current checkpoint",
				UsageText: "rocketpool node claim-ggp [options]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm GGP claim",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return nodeClaimGgp(c)

				},
			},

			{
				Name:      "withdraw-ggp",
				Aliases:   []string{"i"},
				Usage:     "Withdraw GGP staked against the node",
				UsageText: "rocketpool node withdraw-ggp [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of GGP to withdraw (or 'max')",
					},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm GGP withdrawal",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" && c.String("amount") != "max" {
						if _, err := cliutils.ValidatePositiveEthAmount("withdrawal amount", c.String("amount")); err != nil {
							return err
						}
					}

					// Run
					return nodeWithdrawGgp(c)

				},
			},

			{
				Name:      "deposit",
				Aliases:   []string{"d"},
				Usage:     "Make a deposit and create a minipool",
				UsageText: "rocketpool node deposit [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of ETH to deposit (0, 16 or 32)",
					},
					cli.StringFlag{
						Name:  "max-slippage, s",
						Usage: "The maximum acceptable slippage in node commission rate for the deposit (or 'auto')",
					},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm deposit",
					},
					cli.StringFlag{
						Name:  "salt, l",
						Usage: "An optional seed to use when generating the new minipool's address. Use this if you want it to have a custom vanity address.",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" {
						if _, err := cliutils.ValidateDepositEthAmount("deposit amount", c.String("amount")); err != nil {
							return err
						}
					}
					if c.String("max-slippage") != "" && c.String("max-slippage") != "auto" {
						if _, err := cliutils.ValidatePercentage("maximum commission rate slippage", c.String("max-slippage")); err != nil {
							return err
						}
					}
					if c.String("salt") != "" {
						if _, err := cliutils.ValidateBigInt("salt", c.String("salt")); err != nil {
							return err
						}
					}

					// Run
					return nodeDeposit(c)

				},
			},

			{
				Name:      "send",
				Aliases:   []string{"n"},
				Usage:     "Send ETH or tokens from the node account to an address",
				UsageText: "rocketpool node send [options] amount token to",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm token send",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 3); err != nil {
						return err
					}
					amount, err := cliutils.ValidatePositiveEthAmount("send amount", c.Args().Get(0))
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
					return nodeSend(c, amount, token, toAddress)

				},
			},
		},
	})
}
