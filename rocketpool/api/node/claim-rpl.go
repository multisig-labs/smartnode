package node

import (
	"fmt"
	"math/big"

	"github.com/urfave/cli"

	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
)

func canNodeClaimGgp(c *cli.Context) (*api.CanNodeClaimGgpResponse, error) {

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
	response := api.CanNodeClaimGgpResponse{}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Check for rewards
	rewardsAmountWei, err := rewards.GetNodeClaimRewardsAmount(rp, nodeAccount.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting GGP rewards amount: %w", err)
	}
	response.GgpAmount = rewardsAmountWei

	// Don't claim unless the oDAO has claimed first (prevent known issue yet to be patched in smart contracts)
	trustedNodeClaimed, err := rewards.GetTrustedNodeTotalClaimed(rp, nil)
	if err != nil {
		return nil, fmt.Errorf("Error checking if trusted node has already minted GGP: %w", err)
	}
	if trustedNodeClaimed.Cmp(big.NewInt(0)) == 0 {
		response.GgpAmount = big.NewInt(0)
	}

	// Get gas estimate
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}
	gasInfo, err := rewards.EstimateClaimNodeRewardsGas(rp, opts)
	if err != nil {
		return nil, fmt.Errorf("Could not estimate the gas required to claim GGP: %w", err)
	}
	response.GasInfo = gasInfo

	return &response, nil
}

func nodeClaimGgp(c *cli.Context) (*api.NodeClaimGgpResponse, error) {

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
	response := api.NodeClaimGgpResponse{}

	// Get transactor
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}

	// Override the provided pending TX if requested
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}

	// Claim rewards
	hash, err := rewards.ClaimNodeRewards(rp, opts)
	if err != nil {
		return nil, err
	}
	response.TxHash = hash

	// Return response
	return &response, nil

}
