package auction

import (
	"github.com/urfave/cli"

	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage Rocket Pool GGP auctions",
		Subcommands: []cli.Command{

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get GGP auction status",
				UsageText: "rocketpool auction status",
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
				Name:      "lots",
				Aliases:   []string{"l"},
				Usage:     "Get GGP lots for auction",
				UsageText: "rocketpool auction lots",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getLots(c)

				},
			},

			{
				Name:      "create-lot",
				Aliases:   []string{"t"},
				Usage:     "Create a new lot",
				UsageText: "rocketpool auction create-lot",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return createLot(c)

				},
			},

			{
				Name:      "bid-lot",
				Aliases:   []string{"b"},
				Usage:     "Bid on a lot",
				UsageText: "rocketpool auction bid-lot [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "lot, l",
						Usage: "The ID of the lot to bid on",
					},
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of ETH to bid (or 'max')",
					},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm bid",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("lot") != "" {
						if _, err := cliutils.ValidateUint("lot ID", c.String("lot")); err != nil {
							return err
						}
					}
					if c.String("amount") != "" && c.String("amount") != "max" {
						if _, err := cliutils.ValidatePositiveEthAmount("bid amount", c.String("amount")); err != nil {
							return err
						}
					}

					// Run
					return bidOnLot(c)

				},
			},

			{
				Name:      "claim-lot",
				Aliases:   []string{"c"},
				Usage:     "Claim GGP from a lot",
				UsageText: "rocketpool auction claim-lot [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "lot, l",
						Usage: "The lot to claim GGP from (lot ID or 'all')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("lot") != "" && c.String("lot") != "all" {
						if _, err := cliutils.ValidateUint("lot ID", c.String("lot")); err != nil {
							return err
						}
					}

					// Run
					return claimFromLot(c)

				},
			},

			{
				Name:      "recover-lot",
				Aliases:   []string{"r"},
				Usage:     "Recover unclaimed GGP from a lot (returning it to the auction contract)",
				UsageText: "rocketpool auction recover-lot [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "lot, l",
						Usage: "The lot to recover unclaimed GGP from (lot ID or 'all')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("lot") != "" && c.String("lot") != "all" {
						if _, err := cliutils.ValidateUint("lot ID", c.String("lot")); err != nil {
							return err
						}
					}

					// Run
					return recoverGgpFromLot(c)

				},
			},
		},
	})
}
