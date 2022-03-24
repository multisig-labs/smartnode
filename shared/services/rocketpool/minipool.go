package rocketpool

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/smartnode/shared/types/api"
)

// Get minipool status
func (c *Client) MinipoolStatus() (api.MinipoolStatusResponse, error) {
	responseBytes, err := c.callAPI("minipool status")
	if err != nil {
		return api.MinipoolStatusResponse{}, fmt.Errorf("Could not get minipool status: %w", err)
	}
	var response api.MinipoolStatusResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.MinipoolStatusResponse{}, fmt.Errorf("Could not decode minipool status response: %w", err)
	}
	if response.Error != "" {
		return api.MinipoolStatusResponse{}, fmt.Errorf("Could not get minipool status: %s", response.Error)
	}
	for i := 0; i < len(response.Minipools); i++ {
		mp := &response.Minipools[i]
		if mp.Node.DepositBalance == nil {
			mp.Node.DepositBalance = big.NewInt(0)
		}
		if mp.Node.RefundBalance == nil {
			mp.Node.RefundBalance = big.NewInt(0)
		}
		if mp.User.DepositBalance == nil {
			mp.User.DepositBalance = big.NewInt(0)
		}
		// change to AVAX eventually
		if mp.Balances.ETH == nil {
			mp.Balances.ETH = big.NewInt(0)
		}
		if mp.Balances.GGP == nil {
			mp.Balances.GGP = big.NewInt(0)
		}
		if mp.Balances.GGPAVAX == nil {
			mp.Balances.GGPAVAX = big.NewInt(0)
		}
		if mp.Balances.FixedSupplyGGP == nil {
			mp.Balances.FixedSupplyGGP = big.NewInt(0)
		}
		if mp.Validator.Balance == nil {
			mp.Validator.Balance = big.NewInt(0)
		}
		if mp.Validator.NodeBalance == nil {
			mp.Validator.NodeBalance = big.NewInt(0)
		}
	}
	return response, nil
}

// Check whether a minipool is eligible for a refund
func (c *Client) CanRefundMinipool(address common.Address) (api.CanRefundMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-refund %s", address.Hex()))
	if err != nil {
		return api.CanRefundMinipoolResponse{}, fmt.Errorf("Could not get can refund minipool status: %w", err)
	}
	var response api.CanRefundMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanRefundMinipoolResponse{}, fmt.Errorf("Could not decode can refund minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanRefundMinipoolResponse{}, fmt.Errorf("Could not get can refund minipool status: %s", response.Error)
	}
	return response, nil
}

// Refund ETH from a minipool
func (c *Client) RefundMinipool(address common.Address) (api.RefundMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool refund %s", address.Hex()))
	if err != nil {
		return api.RefundMinipoolResponse{}, fmt.Errorf("Could not refund minipool: %w", err)
	}
	var response api.RefundMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.RefundMinipoolResponse{}, fmt.Errorf("Could not decode refund minipool response: %w", err)
	}
	if response.Error != "" {
		return api.RefundMinipoolResponse{}, fmt.Errorf("Could not refund minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool is eligible for staking
func (c *Client) CanStakeMinipool(address common.Address) (api.CanStakeMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-stake %s", address.Hex()))
	if err != nil {
		return api.CanStakeMinipoolResponse{}, fmt.Errorf("Could not get can stake minipool status: %w", err)
	}
	var response api.CanStakeMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanStakeMinipoolResponse{}, fmt.Errorf("Could not decode can stake minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanStakeMinipoolResponse{}, fmt.Errorf("Could not get can stake minipool status: %s", response.Error)
	}
	return response, nil
}

// Stake a minipool
func (c *Client) StakeMinipool(address common.Address) (api.StakeMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool stake %s", address.Hex()))
	if err != nil {
		return api.StakeMinipoolResponse{}, fmt.Errorf("Could not stake minipool: %w", err)
	}
	var response api.StakeMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.StakeMinipoolResponse{}, fmt.Errorf("Could not decode stake minipool response: %w", err)
	}
	if response.Error != "" {
		return api.StakeMinipoolResponse{}, fmt.Errorf("Could not stake minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool can be dissolved
func (c *Client) CanDissolveMinipool(address common.Address) (api.CanDissolveMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-dissolve %s", address.Hex()))
	if err != nil {
		return api.CanDissolveMinipoolResponse{}, fmt.Errorf("Could not get can dissolve minipool status: %w", err)
	}
	var response api.CanDissolveMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanDissolveMinipoolResponse{}, fmt.Errorf("Could not decode can dissolve minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanDissolveMinipoolResponse{}, fmt.Errorf("Could not get can dissolve minipool status: %s", response.Error)
	}
	return response, nil
}

// Dissolve a minipool
func (c *Client) DissolveMinipool(address common.Address) (api.DissolveMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool dissolve %s", address.Hex()))
	if err != nil {
		return api.DissolveMinipoolResponse{}, fmt.Errorf("Could not dissolve minipool: %w", err)
	}
	var response api.DissolveMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.DissolveMinipoolResponse{}, fmt.Errorf("Could not decode dissolve minipool response: %w", err)
	}
	if response.Error != "" {
		return api.DissolveMinipoolResponse{}, fmt.Errorf("Could not dissolve minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool can be exited
func (c *Client) CanExitMinipool(address common.Address) (api.CanExitMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-exit %s", address.Hex()))
	if err != nil {
		return api.CanExitMinipoolResponse{}, fmt.Errorf("Could not get can exit minipool status: %w", err)
	}
	var response api.CanExitMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanExitMinipoolResponse{}, fmt.Errorf("Could not decode can exit minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanExitMinipoolResponse{}, fmt.Errorf("Could not get can exit minipool status: %s", response.Error)
	}
	return response, nil
}

// Exit a minipool
func (c *Client) ExitMinipool(address common.Address) (api.ExitMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool exit %s", address.Hex()))
	if err != nil {
		return api.ExitMinipoolResponse{}, fmt.Errorf("Could not exit minipool: %w", err)
	}
	var response api.ExitMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.ExitMinipoolResponse{}, fmt.Errorf("Could not decode exit minipool response: %w", err)
	}
	if response.Error != "" {
		return api.ExitMinipoolResponse{}, fmt.Errorf("Could not exit minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool can be closed
func (c *Client) CanCloseMinipool(address common.Address) (api.CanCloseMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-close %s", address.Hex()))
	if err != nil {
		return api.CanCloseMinipoolResponse{}, fmt.Errorf("Could not get can close minipool status: %w", err)
	}
	var response api.CanCloseMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanCloseMinipoolResponse{}, fmt.Errorf("Could not decode can close minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanCloseMinipoolResponse{}, fmt.Errorf("Could not get can close minipool status: %s", response.Error)
	}
	return response, nil
}

// Close a minipool
func (c *Client) CloseMinipool(address common.Address) (api.CloseMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool close %s", address.Hex()))
	if err != nil {
		return api.CloseMinipoolResponse{}, fmt.Errorf("Could not close minipool: %w", err)
	}
	var response api.CloseMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CloseMinipoolResponse{}, fmt.Errorf("Could not decode close minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CloseMinipoolResponse{}, fmt.Errorf("Could not close minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool can be finalised
func (c *Client) CanFinaliseMinipool(address common.Address) (api.CanFinaliseMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-finalise %s", address.Hex()))
	if err != nil {
		return api.CanFinaliseMinipoolResponse{}, fmt.Errorf("Could not get can finalise minipool status: %w", err)
	}
	var response api.CanFinaliseMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanFinaliseMinipoolResponse{}, fmt.Errorf("Could not decode can finalise minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanFinaliseMinipoolResponse{}, fmt.Errorf("Could not get can finalise minipool status: %s", response.Error)
	}
	return response, nil
}

// Finalise a minipool
func (c *Client) FinaliseMinipool(address common.Address) (api.FinaliseMinipoolResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool finalise %s", address.Hex()))
	if err != nil {
		return api.FinaliseMinipoolResponse{}, fmt.Errorf("Could not finalise minipool: %w", err)
	}
	var response api.FinaliseMinipoolResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.FinaliseMinipoolResponse{}, fmt.Errorf("Could not decode finalise minipool response: %w", err)
	}
	if response.Error != "" {
		return api.FinaliseMinipoolResponse{}, fmt.Errorf("Could not finalise minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool can have its delegate upgraded
func (c *Client) CanDelegateUpgradeMinipool(address common.Address) (api.CanDelegateUpgradeResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-delegate-upgrade %s", address.Hex()))
	if err != nil {
		return api.CanDelegateUpgradeResponse{}, fmt.Errorf("Could not get can delegate upgrade minipool status: %w", err)
	}
	var response api.CanDelegateUpgradeResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanDelegateUpgradeResponse{}, fmt.Errorf("Could not decode can delegate upgrade minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanDelegateUpgradeResponse{}, fmt.Errorf("Could not get can delegate upgrade minipool status: %s", response.Error)
	}
	return response, nil
}

// Upgrade a minipool delegate
func (c *Client) DelegateUpgradeMinipool(address common.Address) (api.DelegateUpgradeResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool delegate-upgrade %s", address.Hex()))
	if err != nil {
		return api.DelegateUpgradeResponse{}, fmt.Errorf("Could not upgrade delegate for minipool: %w", err)
	}
	var response api.DelegateUpgradeResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.DelegateUpgradeResponse{}, fmt.Errorf("Could not decode upgrade delegate minipool response: %w", err)
	}
	if response.Error != "" {
		return api.DelegateUpgradeResponse{}, fmt.Errorf("Could not upgrade delegate for minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool can have its delegate rolled back
func (c *Client) CanDelegateRollbackMinipool(address common.Address) (api.CanDelegateRollbackResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-delegate-rollback %s", address.Hex()))
	if err != nil {
		return api.CanDelegateRollbackResponse{}, fmt.Errorf("Could not get can delegate rollback minipool status: %w", err)
	}
	var response api.CanDelegateRollbackResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanDelegateRollbackResponse{}, fmt.Errorf("Could not decode can delegate rollback minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanDelegateRollbackResponse{}, fmt.Errorf("Could not get can delegate rollback minipool status: %s", response.Error)
	}
	return response, nil
}

// Rollback a minipool delegate
func (c *Client) DelegateRollbackMinipool(address common.Address) (api.DelegateRollbackResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool delegate-rollback %s", address.Hex()))
	if err != nil {
		return api.DelegateRollbackResponse{}, fmt.Errorf("Could not rollback delegate for minipool: %w", err)
	}
	var response api.DelegateRollbackResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.DelegateRollbackResponse{}, fmt.Errorf("Could not decode rollback delegate minipool response: %w", err)
	}
	if response.Error != "" {
		return api.DelegateRollbackResponse{}, fmt.Errorf("Could not rollback delegate for minipool: %s", response.Error)
	}
	return response, nil
}

// Check whether a minipool can have its auto-upgrade setting changed
func (c *Client) CanSetUseLatestDelegateMinipool(address common.Address, setting bool) (api.CanSetUseLatestDelegateResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool can-set-use-latest-delegate %s %t", address.Hex(), setting))
	if err != nil {
		return api.CanSetUseLatestDelegateResponse{}, fmt.Errorf("Could not get can set use latest delegate for minipool status: %w", err)
	}
	var response api.CanSetUseLatestDelegateResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanSetUseLatestDelegateResponse{}, fmt.Errorf("Could not decode can set use latest delegate for minipool response: %w", err)
	}
	if response.Error != "" {
		return api.CanSetUseLatestDelegateResponse{}, fmt.Errorf("Could not get can set use latest delegate for minipool status: %s", response.Error)
	}
	return response, nil
}

// Change a minipool's auto-upgrade setting
func (c *Client) SetUseLatestDelegateMinipool(address common.Address, setting bool) (api.SetUseLatestDelegateResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool set-use-latest-delegate %s %t", address.Hex(), setting))
	if err != nil {
		return api.SetUseLatestDelegateResponse{}, fmt.Errorf("Could not set use latest delegate for minipool: %w", err)
	}
	var response api.SetUseLatestDelegateResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.SetUseLatestDelegateResponse{}, fmt.Errorf("Could not decode set use latest delegate for minipool response: %w", err)
	}
	if response.Error != "" {
		return api.SetUseLatestDelegateResponse{}, fmt.Errorf("Could not set use latest delegate for minipool: %s", response.Error)
	}
	return response, nil
}

// Get the artifacts necessary for vanity address searching
func (c *Client) GetVanityArtifacts(depositAmount *big.Int, nodeAddress string) (api.GetVanityArtifactsResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("minipool get-vanity-artifacts %s %s", depositAmount.String(), nodeAddress))
	if err != nil {
		return api.GetVanityArtifactsResponse{}, fmt.Errorf("Could not get vanity artifacts: %w", err)
	}
	var response api.GetVanityArtifactsResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.GetVanityArtifactsResponse{}, fmt.Errorf("Could not decode get vanity artifacts response: %w", err)
	}
	if response.Error != "" {
		return api.GetVanityArtifactsResponse{}, fmt.Errorf("Could not get vanity artifacts: %s", response.Error)
	}
	return response, nil
}
