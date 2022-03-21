package faucet

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
		Usage:   "Access the legacy GGP faucet",
		Subcommands: []cli.Command{

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get the faucet's status",
				UsageText: "rocketpool api faucet status",
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
				Name:      "can-withdraw-ggp",
				Usage:     "Check whether the node can withdraw legacy GGP from the faucet",
				UsageText: "rocketpool api faucet can-withdraw-ggp",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(canWithdrawGgp(c))
					return nil

				},
			},
			{
				Name:      "withdraw-ggp",
				Aliases:   []string{"w"},
				Usage:     "Withdraw legacy GGP from the faucet",
				UsageText: "rocketpool api faucet withdraw-ggp",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					api.PrintResponse(withdrawGgp(c))
					return nil

				},
			},
		},
	})
}
