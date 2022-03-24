package node

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
)

func canNodeSwapGgp(c *cli.Context, amountWei *big.Int) (*api.CanNodeSwapGgpResponse, error) {

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
	response := api.CanNodeSwapGgpResponse{}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Check node fixed-supply GGP balance
	fixedSupplyGgpBalance, err := tokens.GetFixedSupplyGGPBalance(rp, nodeAccount.Address, nil)
	if err != nil {
		return nil, err
	}
	response.InsufficientBalance = (amountWei.Cmp(fixedSupplyGgpBalance) > 0)

	// Get gas estimates
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	gasInfo, err := tokens.EstimateSwapFixedSupplyGGPForGGPGas(rp, amountWei, opts)
	if err != nil {
		return nil, err
	}
	response.GasInfo = gasInfo

	// Update & return response
	response.CanSwap = !response.InsufficientBalance
	return &response, nil

}

func allowanceFsGgp(c *cli.Context) (*api.NodeSwapGgpAllowanceResponse, error) {

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
	response := api.NodeSwapGgpAllowanceResponse{}

	// Get new GGP contract address
	gogoTokenGGPAddress, err := rp.GetAddress("gogoTokenGGP")
	if err != nil {
		return nil, err
	}

	// Get node account
	account, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Get node's FSGGP allowance
	allowance, err := tokens.GetFixedSupplyGGPAllowance(rp, account.Address, *gogoTokenGGPAddress, nil)
	if err != nil {
		return nil, err
	}

	response.Allowance = allowance

	return &response, nil
}

func getSwapApprovalGas(c *cli.Context, amountWei *big.Int) (*api.NodeSwapGgpApproveGasResponse, error) {
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
	response := api.NodeSwapGgpApproveGasResponse{}

	// Get GGP contract address
	gogoTokenGGPAddress, err := rp.GetAddress("gogoTokenGGP")
	if err != nil {
		return nil, err
	}

	// Get gas estimates
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	gasInfo, err := tokens.EstimateApproveFixedSupplyGGPGas(rp, *gogoTokenGGPAddress, amountWei, opts)
	if err != nil {
		return nil, err
	}
	response.GasInfo = gasInfo
	return &response, nil
}

func approveFsGgp(c *cli.Context, amountWei *big.Int) (*api.NodeSwapGgpApproveResponse, error) {

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
	response := api.NodeSwapGgpApproveResponse{}

	// Get GGP contract address
	gogoTokenGGPAddress, err := rp.GetAddress("gogoTokenGGP")
	if err != nil {
		return nil, err
	}

	// Approve fixed-supply GGP allowance
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}
	if hash, err := tokens.ApproveFixedSupplyGGP(rp, *gogoTokenGGPAddress, amountWei, opts); err != nil {
		return nil, err
	} else {
		response.ApproveTxHash = hash
	}

	// Return response
	return &response, nil

}

func waitForApprovalAndSwapFsGgp(c *cli.Context, amountWei *big.Int, hash common.Hash) (*api.NodeSwapGgpSwapResponse, error) {

	// Get services
	if err := services.RequireNodeWallet(c); err != nil {
		return nil, err
	}
	if err := services.RequireRocketStorage(c); err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Wait for the fixed-supply GGP approval TX to successfully get mined
	_, err = utils.WaitForTransaction(rp.Client, hash)
	if err != nil {
		return nil, err
	}

	return swapGgp(c, amountWei)

}

func swapGgp(c *cli.Context, amountWei *big.Int) (*api.NodeSwapGgpSwapResponse, error) {

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
	response := api.NodeSwapGgpSwapResponse{}

	// Swap fixed-supply GGP for GGP
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}
	if hash, err := tokens.SwapFixedSupplyGGPForGGP(rp, amountWei, opts); err != nil {
		return nil, err
	} else {
		response.SwapTxHash = hash
	}

	// Return response
	return &response, nil

}
