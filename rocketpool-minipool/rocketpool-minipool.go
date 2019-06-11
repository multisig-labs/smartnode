package main

import (
    "log"
    "os"

    "github.com/ethereum/go-ethereum/common"
    "gopkg.in/urfave/cli.v1"

    "github.com/rocket-pool/smartnode/rocketpool-minipool/minipool"
    "github.com/rocket-pool/smartnode/shared/services"
    cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"

    "fmt"
)


// Run application
func main() {

    // Initialise application
    app := cli.NewApp()

    // Set application info
    app.Name = "rocketpool-minipool"
    app.Usage = "Rocket Pool minipool activity daemon"
    app.Version = "0.0.1"
    app.Authors = []cli.Author{
        cli.Author{
            Name:  "Darren Langley",
            Email: "darren@rocketpool.net",
        },
        cli.Author{
            Name:  "David Rugendyke",
            Email: "david@rocketpool.net",
        },
        cli.Author{
            Name:  "Jake Pospischil",
            Email: "jake@rocketpool.net",
        },
    }
    app.Copyright = "(c) 2019 Rocket Pool Pty Ltd"

    // Configure application
    cliutils.Configure(app)

    // Set application action
    app.Action = func(c *cli.Context) error {

        // Check argument count
        if len(c.Args()) != 1 {
            return cli.NewExitError("USAGE:" + "\n   " + "rocketpool-minipool address", 1)
        }

        // Get & validate arguments
        address := c.Args().Get(0)
        if !common.IsHexAddress(address) {
            return cli.NewExitError("Invalid minipool address", 1)
        }

        // Run process
        return run(c, address)

    }

    // Run application
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }

}


// Run process
func run(c *cli.Context, address string) error {

    // Initialise services
    p, err := services.NewProvider(c, services.ProviderOpts{
        Client: true,
        ClientSync: true,
        Publisher: true,
        Beacon: true,
        LoadContracts: []string{},
        LoadAbis: []string{"rocketMinipool"},
    })
    if err != nil {
        return err
    }

    // Start minipool processes
    go minipool.StartActivityProcess(p)
    go minipool.StartWithdrawalProcess(p)

    // Start services
    p.Beacon.Connect()

    // Block thread
    select {}

}
