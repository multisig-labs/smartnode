package node

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	tndao "github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	tnsettings "github.com/rocket-pool/rocketpool-go/settings/trustednode"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
	"github.com/rocket-pool/smartnode/shared/utils/validator"
)


type minipoolCreated struct {
    Minipool common.Address
    Node common.Address
    Time *big.Int
}


func canNodeDeposit(c *cli.Context, amountWei *big.Int, minNodeFee float64, salt *big.Int) (*api.CanNodeDepositResponse, error) {

    // Get services
    if err := services.RequireNodeRegistered(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    ec, err := services.GetEthClient(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }
    bc, err := services.GetBeaconClient(c)
    if err != nil { return nil, err }

    // Get eth2 config
    eth2Config, err := bc.GetEth2Config()
    if err != nil {
        return nil, err
    }

    // Response
    response := api.CanNodeDepositResponse{}

    // Check if amount is zero
    amountIsZero := (amountWei.Cmp(big.NewInt(0)) == 0)

    // Get node account
    nodeAccount, err := w.GetNodeAccount()
    if err != nil {
        return nil, err
    }

    // Adjust the salt
    if salt.Cmp(big.NewInt(0)) == 0 {
        nonce, err := ec.NonceAt(context.Background(), nodeAccount.Address, nil)
        if err != nil {
            return nil, err
        }
        salt.SetUint64(nonce)
    }

    // Data
    var wg1 errgroup.Group
    var isTrusted bool
    var minipoolCount uint64
    var minipoolLimit uint64

    // Check node balance
    wg1.Go(func() error {
        ethBalanceWei, err := ec.BalanceAt(context.Background(), nodeAccount.Address, nil)
        if err == nil {
            response.InsufficientBalance = (amountWei.Cmp(ethBalanceWei) > 0)
        }
        return err
    })

    // Check node deposits are enabled
    wg1.Go(func() error {
        depositEnabled, err := protocol.GetNodeDepositEnabled(rp, nil)
        if err == nil {
            response.DepositDisabled = !depositEnabled
        }
        return err
    })

    // Get trusted status
    wg1.Go(func() error {
        var err error
        isTrusted, err = tndao.GetMemberExists(rp, nodeAccount.Address, nil)
        return err
    })

    // Get node staking information
    wg1.Go(func() error {
        var err error
        minipoolCount, err = minipool.GetNodeMinipoolCount(rp, nodeAccount.Address, nil)
        return err
    })
    wg1.Go(func() error {
        var err error
        minipoolLimit, err = node.GetNodeMinipoolLimit(rp, nodeAccount.Address, nil)
        return err
    })

    // Get consensus status
    wg1.Go(func() error {
        var err error
        inConsensus, err := network.InConsensus(rp, nil)
        response.InConsensus = inConsensus
        return err
    })

    // Get gas estimate
    wg1.Go(func() error {
        opts, err := w.GetNodeAccountTransactor()
        if err != nil { 
            return err 
        }
        opts.Value = amountWei

        // Get the deposit type
        depositType, err := node.GetDepositType(rp, amountWei, nil)
        if err != nil {
            return err
        }

        // Get the next validator key
        validatorKey, err := w.GetNextValidatorKey()
        if err != nil {
            return err
        }

        // Get the next minipool address and withdrawal credentials
        minipoolAddress, err := utils.GenerateAddress(rp, nodeAccount.Address, depositType, salt, nil)
        if err != nil {
            return err
        }
        withdrawalCredentials := utils.GetWithdrawalCredentials(minipoolAddress)

        // Get validator deposit data and associated parameters
        depositData, depositDataRoot, err := validator.GetDepositData(validatorKey, withdrawalCredentials, eth2Config)
        if err != nil {
            return err
        }
        pubKey := rptypes.BytesToValidatorPubkey(depositData.PublicKey)
        signature := rptypes.BytesToValidatorSignature(depositData.Signature)

        // Run the deposit gas estimator
        gasInfo, err := node.EstimateDepositGas(rp, minNodeFee, pubKey, signature, depositDataRoot, salt, minipoolAddress, opts)
        if err == nil {
            response.GasInfo = gasInfo
        }
        return err
    })

    // Wait for data
    if err := wg1.Wait(); err != nil {
        return nil, err
    }

    // Check data
    response.InsufficientRplStake = (minipoolCount >= minipoolLimit)
    response.InvalidAmount = (!isTrusted && amountIsZero)

    // Check oracle node unbonded minipool limit
    if isTrusted && amountIsZero {

        // Data
        var wg2 errgroup.Group
        var unbondedMinipoolCount uint64
        var unbondedMinipoolsMax uint64

        // Get unbonded minipool details
        wg2.Go(func() error {
            var err error
            unbondedMinipoolCount, err = tndao.GetMemberUnbondedValidatorCount(rp, nodeAccount.Address, nil)
            return err
        })
        wg2.Go(func() error {
            var err error
            unbondedMinipoolsMax, err = tnsettings.GetMinipoolUnbondedMax(rp, nil)
            return err
        })

        // Wait for data
        if err := wg2.Wait(); err != nil {
            return nil, err
        }

        // Check unbonded minipool limit
        response.UnbondedMinipoolsAtMax = (unbondedMinipoolCount >= unbondedMinipoolsMax)

    }

    // Update & return response
    response.CanDeposit = !(response.InsufficientBalance || response.InsufficientRplStake || response.InvalidAmount || response.UnbondedMinipoolsAtMax || response.DepositDisabled || !response.InConsensus)
    return &response, nil

}


func nodeDeposit(c *cli.Context, amountWei *big.Int, minNodeFee float64, salt *big.Int) (*api.NodeDepositResponse, error) {

    // Get services
    if err := services.RequireNodeRegistered(c); err != nil { return nil, err }
    w, err := services.GetWallet(c)
    if err != nil { return nil, err }
    ec, err := services.GetEthClient(c)
    if err != nil { return nil, err }
    rp, err := services.GetRocketPool(c)
    if err != nil { return nil, err }
    bc, err := services.GetBeaconClient(c)
    if err != nil { return nil, err }

    // Get eth2 config
    eth2Config, err := bc.GetEth2Config()
    if err != nil {
        return nil, err
    }

    // Get node account
    nodeAccount, err := w.GetNodeAccount()
    if err != nil {
        return nil, err
    }

    // Response
    response := api.NodeDepositResponse{}

    // Adjust the salt
    if salt.Cmp(big.NewInt(0)) == 0 {
        nonce, err := ec.NonceAt(context.Background(), nodeAccount.Address, nil)
        if err != nil {
            return nil, err
        }
        salt.SetUint64(nonce)
    }

    // Make sure ETH2 is on the correct chain
    depositContractInfo, err := getDepositContractInfo(c)
    if err != nil {
        return nil, err
    }
    if depositContractInfo.RPNetwork != depositContractInfo.BeaconNetwork ||
       depositContractInfo.RPDepositContract != depositContractInfo.BeaconDepositContract {
            return nil, fmt.Errorf("Beacon network mismatch! Expected %s on chain %d, but beacon is using %s on chain %d.",
                            depositContractInfo.RPDepositContract.Hex(),
                            depositContractInfo.RPNetwork,
                            depositContractInfo.BeaconDepositContract.Hex(),
                            depositContractInfo.BeaconNetwork)
    }

    // Get transactor
    opts, err := w.GetNodeAccountTransactor()
    if err != nil {
        return nil, err
    }
    opts.Value = amountWei

    // Get the deposit type
    depositType, err := node.GetDepositType(rp, amountWei, nil)
    if err != nil {
        return nil, err
    }

    // Create and save a new validator key
    validatorKey, err := w.CreateValidatorKey()
    if err != nil {
        return nil, err
    }

    // Get the next minipool address and withdrawal credentials
    minipoolAddress, err := utils.GenerateAddress(rp, nodeAccount.Address, depositType, salt, nil)
    if err != nil {
        return nil, err
    }
    withdrawalCredentials := utils.GetWithdrawalCredentials(minipoolAddress)

    // Get validator deposit data and associated parameters
    depositData, depositDataRoot, err := validator.GetDepositData(validatorKey, withdrawalCredentials, eth2Config)
    if err != nil {
        return nil, err
    }
    pubKey := rptypes.BytesToValidatorPubkey(depositData.PublicKey)
    signature := rptypes.BytesToValidatorSignature(depositData.Signature)

    // Override the provided pending TX if requested 
    err = eth1.CheckForNonceOverride(c, opts)
    if err != nil {
        return nil, fmt.Errorf("Error checking for nonce override: %w", err)
    }

    // Deposit
    hash, err := node.Deposit(rp, minNodeFee, pubKey, signature, depositDataRoot, salt, minipoolAddress, opts)
    if err != nil {
        return nil, err
    }
    response.TxHash = hash
    response.MinipoolAddress = minipoolAddress

    // Return response
    return &response, nil

}

