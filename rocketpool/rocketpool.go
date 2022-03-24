package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/rocketpool/api"
	"github.com/rocket-pool/smartnode/rocketpool/node"
	"github.com/rocket-pool/smartnode/rocketpool/watchtower"
	"github.com/rocket-pool/smartnode/shared"
	apiutils "github.com/rocket-pool/smartnode/shared/utils/api"
)

// Run
func main() {

	// Initialise application
	app := cli.NewApp()

	// Set application info
	app.Name = "rocketpool"
	app.Usage = "Rocket Pool service"
	app.Version = shared.RocketPoolVersion
	app.Authors = []cli.Author{
		{
			Name:  "David Rugendyke",
			Email: "david@rocketpool.net",
		},
		{
			Name:  "Jake Pospischil",
			Email: "jake@rocketpool.net",
		},
		{
			Name:  "Joe Clapis",
			Email: "joe@rocketpool.net",
		},
		{
			Name:  "Kane Wallmann",
			Email: "kane@rocketpool.net",
		},
	}
	app.Copyright = "(c) 2021 Rocket Pool Pty Ltd"

	// Set application flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Rocket Pool service global config absolute `path`",
			Value: "/.rocketpool/config.yml",
		},
		cli.StringFlag{
			Name:  "settings, s",
			Usage: "Rocket Pool service user config absolute `path`",
			Value: "/.rocketpool/settings.yml",
		},
		cli.StringFlag{
			Name:  "storageAddress, a",
			Usage: "Rocket Pool storage contract `address`",
		},
		cli.StringFlag{
			Name:  "oneInchOracleAddress, o",
			Usage: "1inch exchange oracle contract `address`",
		},
		cli.StringFlag{
			Name:  "ggpTokenAddress, t",
			Usage: "GGP token contract `address`",
		},
		cli.StringFlag{
			Name:  "ggpFaucetAddress, f",
			Usage: "Rocket Pool GGP token faucet `address`",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "Rocket Pool wallet password file absolute `path`",
		},
		cli.StringFlag{
			Name:  "wallet, w",
			Usage: "Rocket Pool wallet file absolute `path`",
		},
		cli.StringFlag{
			Name:  "validatorKeychain, k",
			Usage: "Rocket Pool validator keychain absolute `path`",
		},
		cli.StringFlag{
			Name:  "eth1Provider, e",
			Usage: "Eth 1.0 provider `address`",
		},
		cli.StringFlag{
			Name:  "eth2Provider, b",
			Usage: "Eth 2.0 provider `address`",
		},
		cli.Float64Flag{
			Name:  "maxFee",
			Usage: "Desired max fee in gwei",
		},
		cli.Float64Flag{
			Name:  "maxPrioFee",
			Usage: "Desired max priority fee in gwei",
		},
		cli.Uint64Flag{
			Name:  "gasLimit, l",
			Usage: "Desired gas limit",
		},
		cli.StringFlag{
			Name:  "nonce",
			Usage: "Use this flag to explicitly specify the nonce that this transaction should use, so it can override an existing 'stuck' transaction",
		},
		cli.StringFlag{
			Name:  "metricsAddress, m",
			Usage: "Address to serve metrics on if enabled",
			Value: "0.0.0.0",
		},
		cli.UintFlag{
			Name:  "metricsPort, r",
			Usage: "Port to serve metrics on if enabled",
			Value: 9102,
		},
	}

	// Register commands
	api.RegisterCommands(app, "api", []string{"a"})
	node.RegisterCommands(app, "node", []string{"n"})
	watchtower.RegisterCommands(app, "watchtower", []string{"w"})

	// Get command being run
	var commandName string
	app.Before = func(c *cli.Context) error {
		commandName = c.Args().First()
		return nil
	}

	// Run application
	if err := app.Run(os.Args); err != nil {
		if commandName == "api" {
			apiutils.PrintErrorResponse(err)
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

}
