package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
)

func canNodeStakeGgp(c *cli.Context, amountWei *big.Int) (*api.CanNodeStakeGgpResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.CanNodeStakeGgpResponse{}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Check GGP balance
	ggpBalance, err := tokens.GetGGPBalance(rp, nodeAccount.Address, nil)
	if err != nil {
		return nil, err
	}
	response.InsufficientBalance = (amountWei.Cmp(ggpBalance) > 0)

	// Check network consensus
	inConsensus, err := network.InConsensus(rp, nil)
	if err != nil {
		return nil, err
	}
	response.InConsensus = inConsensus

	// Get gas estimates
	opts, err := w.GetNodeAccountTransactor()
	opts.GasPrice = big.NewInt(225000000000)
	opts.GasLimit = uint64(225000000000)
	opts.GasFeeCap = big.NewInt(225000000000)
	opts.GasTipCap = big.NewInt(225000000000)
	if err != nil {
		return nil, err
	}
	gasInfo, err := node.EstimateStakeGas(rp, amountWei, opts)
	if err != nil {
		return nil, err
	}
	response.GasInfo = gasInfo

	// Update & return response
	response.CanStake = !(response.InsufficientBalance || !response.InConsensus)
	return &response, nil

}

func getStakeApprovalGas(c *cli.Context, amountWei *big.Int) (*api.NodeStakeGgpApproveGasResponse, error) {
	// Get services
	if err := services.RequireNodeWallet(c); err != nil {
		return nil, err
	}
	if err := services.RequireRocketStorage(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NodeStakeGgpApproveGasResponse{}

	// Get staking contract address
	rocketNodeStakingAddress, err := rp.GetAddress("rocketNodeStaking")
	if err != nil {
		return nil, err
	}

	// Get gas estimates
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	gasInfo, err := tokens.EstimateApproveGGPGas(rp, *rocketNodeStakingAddress, amountWei, opts)
	if err != nil {
		return nil, err
	}
	response.GasInfo = gasInfo
	return &response, nil
}

func allowanceGgp(c *cli.Context) (*api.NodeStakeGgpAllowanceResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NodeStakeGgpAllowanceResponse{}

	// Get staking contract address
	rocketNodeStakingAddress, err := rp.GetAddress("rocketNodeStaking")
	if err != nil {
		return nil, err
	}

	// Get node account
	account, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Get node's GGP allowance
	allowance, err := tokens.GetGGPAllowance(rp, account.Address, *rocketNodeStakingAddress, nil)
	if err != nil {
		return nil, err
	}

	response.Allowance = allowance

	return &response, nil
}

func approveGgp(c *cli.Context, amountWei *big.Int) (*api.NodeStakeGgpApproveResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NodeStakeGgpApproveResponse{}

	// Get staking contract address
	rocketNodeStakingAddress, err := rp.GetAddress("rocketNodeStaking")
	if err != nil {
		return nil, err
	}

	// Approve GGP allowance
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}
	if hash, err := tokens.ApproveGGP(rp, *rocketNodeStakingAddress, amountWei, opts); err != nil {
		return nil, err
	} else {
		response.ApproveTxHash = hash
	}

	// Return response
	return &response, nil

}

func waitForApprovalAndStakeGgp(c *cli.Context, amountWei *big.Int, hash common.Hash) (*api.NodeStakeGgpStakeResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Wait for the GGP approval TX to successfully get mined
	_, err = utils.WaitForTransaction(rp.Client, hash)
	if err != nil {
		return nil, err
	}

	// Perform the stake
	return stakeGgp(c, amountWei)

}

func stakeGgp(c *cli.Context, amountWei *big.Int) (*api.NodeStakeGgpStakeResponse, error) {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.NodeStakeGgpStakeResponse{}

	// Stake GGP
	opts, err := w.GetNodeAccountTransactor()
	fmt.Println(opts.From.String())
	if err != nil {
		return nil, err
	}
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}
	if hash, err := node.StakeGGP(rp, amountWei, opts); err != nil {
		return nil, err
	} else {
		response.StakeTxHash = hash
	}

	// Return response
	return &response, nil

}
