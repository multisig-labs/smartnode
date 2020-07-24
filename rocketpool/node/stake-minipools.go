package node

import (
    "fmt"
    "log"
    "time"

    "github.com/ethereum/go-ethereum/common"
    "github.com/rocket-pool/rocketpool-go/minipool"
    "github.com/rocket-pool/rocketpool-go/rocketpool"
    "github.com/rocket-pool/rocketpool-go/types"
    "github.com/urfave/cli"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/services/accounts"
)


// Settings
var checkMinipoolsInterval, _ = time.ParseDuration("1m")


// Start stake prelaunch minipools task
func startStakePrelaunchMinipools(c *cli.Context) error {

    // Get services
    if err := services.WaitNodeRegistered(c, true); err != nil { return err }
    am, err := services.GetAccountManager(c)
    if err != nil { return err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return err }

    // Stake prelaunch minipools at interval
    go (func() {
        for {
            if err := stakePrelaunchMinipools(c, am, rp); err != nil {
                log.Println(err)
            }
            time.Sleep(checkMinipoolsInterval)
        }
    })()

    // Return
    return nil

}


// Stake prelaunch minipools
func stakePrelaunchMinipools(c *cli.Context, am *accounts.AccountManager, rp *rocketpool.RocketPool) error {

    // Wait for eth client to sync
    if err := services.WaitClientSynced(c, true); err != nil {
        return err
    }

    // Get node account
    nodeAccount, err := am.GetNodeAccount()
    if err != nil {
        return err
    }

    // Get prelaunch minipools
    minipools, err := getPrelaunchMinipools(rp, nodeAccount.Address)
    if err != nil {
        return err
    }
    if len(minipools) == 0 {
        return nil
    }

    // Stake minipools
    log.Printf("%d minipools are ready for staking...\n", len(minipools))
    for _, mp := range minipools {
        log.Printf("Staking minipool %s...\n", mp.Address.Hex())
        if err := stakeMinipool(rp, mp); err != nil {
            return fmt.Errorf("Could not stake minipool %s: %w", mp.Address.Hex(), err)
        }
        log.Printf("Successfully staked minipool %s\n", mp.Address.Hex())
    }

    // Return
    return nil

}


// Get prelaunch minipools
func getPrelaunchMinipools(rp *rocketpool.RocketPool, nodeAddress common.Address) ([]*minipool.Minipool, error) {

    // Get node minipool addresses
    addresses, err := minipool.GetNodeMinipoolAddresses(rp, nodeAddress)
    if err != nil {
        return []*minipool.Minipool{}, err
    }

    // Create minipool contracts
    minipools := make([]*minipool.Minipool, len(addresses))
    for mi, address := range addresses {
        mp, err := minipool.NewMinipool(rp, address)
        if err != nil {
            return []*minipool.Minipool{}, err
        }
        minipools[mi] = mp
    }

    // Data
    var wg errgroup.Group
    statuses := make([]types.MinipoolStatus, len(minipools))

    // Load minipool statuses
    for mi, mp := range minipools {
        mi, mp := mi, mp
        wg.Go(func() error {
            status, err := mp.GetStatus()
            if err == nil { statuses[mi] = status }
            return err
        })
    }

    // Wait for data
    if err := wg.Wait(); err != nil {
        return []*minipool.Minipool{}, err
    }

    // Filter minipools by status
    prelaunchMinipools := []*minipool.Minipool{}
    for mi, mp := range minipools {
        if statuses[mi] == types.Prelaunch {
            prelaunchMinipools = append(prelaunchMinipools, mp)
        }
    }

    // Return
    return prelaunchMinipools, nil

}


// Stake a minipool
func stakeMinipool(rp *rocketpool.RocketPool, mp *minipool.Minipool) error {

    // TODO: implement
    log.Println("Minipool staking not implemented...")
    return nil

}
