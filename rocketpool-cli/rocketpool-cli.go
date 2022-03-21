package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/rocketpool-cli/auction"
	"github.com/rocket-pool/smartnode/rocketpool-cli/faucet"
	"github.com/rocket-pool/smartnode/rocketpool-cli/minipool"
	"github.com/rocket-pool/smartnode/rocketpool-cli/network"
	"github.com/rocket-pool/smartnode/rocketpool-cli/node"
	"github.com/rocket-pool/smartnode/rocketpool-cli/odao"
	"github.com/rocket-pool/smartnode/rocketpool-cli/queue"
	"github.com/rocket-pool/smartnode/rocketpool-cli/service"
	"github.com/rocket-pool/smartnode/rocketpool-cli/wallet"
	"github.com/rocket-pool/smartnode/shared"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

// Run
func main() {

	// Add logo to application help template
	cli.AppHelpTemplate = fmt.Sprintf(`
______           _        _    ______           _ 
| ___ \         | |      | |   | ___ \         | |
| |_/ /___   ___| | _____| |_  | |_/ /__   ___ | |
|    // _ \ / __| |/ / _ \ __| |  __/ _ \ / _ \| |
| |\ \ (_) | (__|   <  __/ |_  | | | (_) | (_) | |
\_| \_\___/ \___|_|\_\___|\__| \_|  \___/ \___/|_|

%s`, cli.AppHelpTemplate)

	// Initialise application
	app := cli.NewApp()

	// Set application info
	app.Name = "rocketpool"
	app.Usage = "Rocket Pool CLI"
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
		cli.BoolFlag{
			Name:  "allow-root, r",
			Usage: "Allow rocketpool to be run as the root user",
		},
		cli.StringFlag{
			Name:  "config-path, c",
			Usage: "Rocket Pool config asset `path`",
			Value: "~/.rocketpool",
		},
		cli.StringFlag{
			Name:  "daemon-path, d",
			Usage: "Interact with a Rocket Pool service daemon at a `path` on the host OS, running outside of docker",
		},
		cli.StringFlag{
			Name:  "host, o",
			Usage: "DEPRECATED - Smart node SSH host `address`",
		},
		cli.StringFlag{
			Name:  "user, u",
			Usage: "DEPRECATED - Smart node SSH user `name`",
		},
		cli.StringFlag{
			Name:  "key, k",
			Usage: "DEPRECATED - Smart node SSH key `file`",
		},
		cli.StringFlag{
			Name:  "passphrase, p",
			Usage: "DEPRECATED - Smart node SSH key passphrase `file`",
		},
		cli.StringFlag{
			Name:  "known-hosts, n",
			Usage: "DEPRECATED - Smart node SSH known_hosts `file` (default: current user's ~/.ssh/known_hosts)",
		},
		cli.StringFlag{
			Name:  "gasPrice, g",
			Usage: "OBSOLETE - No longer used, please use --maxFee and --maxPrioFee instead",
		},
		cli.Float64Flag{
			Name:  "maxFee, f",
			Usage: "The max fee (including the priority fee) you want a transaction to cost, in gwei",
		},
		cli.Float64Flag{
			Name:  "maxPrioFee, i",
			Usage: "The max priority fee you want a transaction to use, in gwei",
		},
		cli.Uint64Flag{
			Name:  "gasLimit, l",
			Usage: "Desired gas limit",
		},
		cli.StringFlag{
			Name:  "nonce",
			Usage: "Use this flag to explicitly specify the nonce that this transaction should use, so it can override an existing 'stuck' transaction",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug printing of API commands",
		},
		cli.BoolFlag{
			Name: "secure-session, s",
			Usage: "Some commands may print sensitive information to your terminal. " +
				"Use this flag when nobody can see your screen to allow sensitive data to be printed without prompting",
		},
	}

	// Register commands
	auction.RegisterCommands(app, "auction", []string{"a"})

	// Get the config path from the arguments (or use the default)
	configPath := "~/.rocketpool"
	for index, arg := range os.Args {
		if arg == "-c" || arg == "--config-path" {
			if len(os.Args)-1 == index {
				fmt.Fprintf(os.Stderr, "Expected config path after %s but none was given.\n", arg)
				os.Exit(1)
			}
			configPath = os.Args[index+1]
		}
	}

	// Get and parse the config file
	configFile := fmt.Sprintf("%s/%s", configPath, rocketpool.GlobalConfigFile)
	expandedPath, err := homedir.Expand(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the global config file path: %s\n", err.Error())
		os.Exit(1)
	}
	// Stop if the config file doesn't exist yet
	_, err = os.Stat(expandedPath)
	if !os.IsNotExist(err) {
		configBytes, err := ioutil.ReadFile(expandedPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load the global config file: %s\n", err.Error())
			os.Exit(1)
		}
		cfg, err := config.Parse(configBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse the global config file: %s\n", err.Error())
			os.Exit(1)
		}

		// Add the faucet if we're on a testnet and it has a contract address
		if cfg.Rocketpool.GGPFaucetAddress != "" {
			faucet.RegisterCommands(app, "faucet", []string{"f"})
		}
	}

	minipool.RegisterCommands(app, "minipool", []string{"m"})
	network.RegisterCommands(app, "network", []string{"e"})
	node.RegisterCommands(app, "node", []string{"n"})
	odao.RegisterCommands(app, "odao", []string{"o"})
	queue.RegisterCommands(app, "queue", []string{"q"})
	service.RegisterCommands(app, "service", []string{"s"})
	wallet.RegisterCommands(app, "wallet", []string{"w"})

	app.Before = func(c *cli.Context) error {
		// Check user ID
		if os.Getuid() == 0 && !c.GlobalBool("allow-root") {
			fmt.Fprintln(os.Stderr, "rocketpool should not be run as root. Please try again without 'sudo'.")
			fmt.Fprintln(os.Stderr, "If you want to run rocketpool as root anyway, use the '--allow-root' option to override this warning.")
			os.Exit(1)
		}

		// Check for deprecated flags
		if c.String("gasPrice") != "" {
			fmt.Fprintln(os.Stderr, "The `gasPrice` flag is deprecated - please use `--maxFee` and optionally `--maxPrioFee` instead.")
			os.Exit(1)
		}

		return nil
	}

	// Run application
	fmt.Println("")
	if err := app.Run(os.Args); err != nil {
		cliutils.PrettyPrintError(err)
	}
	fmt.Println("")

}
