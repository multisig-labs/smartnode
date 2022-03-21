package rocketpool

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/smartnode/shared/types/api"
)

// Get node status
func (c *Client) NodeStatus() (api.NodeStatusResponse, error) {
	responseBytes, err := c.callAPI("node status")
	if err != nil {
		return api.NodeStatusResponse{}, fmt.Errorf("Could not get node status: %w", err)
	}
	var response api.NodeStatusResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeStatusResponse{}, fmt.Errorf("Could not decode node status response: %w", err)
	}
	if response.Error != "" {
		return api.NodeStatusResponse{}, fmt.Errorf("Could not get node status: %s", response.Error)
	}
	if response.GgpStake == nil {
		response.GgpStake = big.NewInt(0)
	}
	if response.EffectiveGgpStake == nil {
		response.EffectiveGgpStake = big.NewInt(0)
	}
	if response.MinimumGgpStake == nil {
		response.MinimumGgpStake = big.NewInt(0)
	}
	if response.AccountBalances.ETH == nil {
		response.AccountBalances.ETH = big.NewInt(0)
	}
	if response.AccountBalances.GGP == nil {
		response.AccountBalances.GGP = big.NewInt(0)
	}
	if response.AccountBalances.GGPAVAX == nil {
		response.AccountBalances.GGPAVAX = big.NewInt(0)
	}
	if response.AccountBalances.FixedSupplyGGP == nil {
		response.AccountBalances.FixedSupplyGGP = big.NewInt(0)
	}
	if response.WithdrawalBalances.ETH == nil {
		response.WithdrawalBalances.ETH = big.NewInt(0)
	}
	if response.WithdrawalBalances.GGP == nil {
		response.WithdrawalBalances.GGP = big.NewInt(0)
	}
	if response.WithdrawalBalances.GGPAVAX == nil {
		response.WithdrawalBalances.GGPAVAX = big.NewInt(0)
	}
	if response.WithdrawalBalances.FixedSupplyGGP == nil {
		response.WithdrawalBalances.FixedSupplyGGP = big.NewInt(0)
	}
	return response, nil
}

// Check whether the node can be registered
func (c *Client) CanRegisterNode(timezoneLocation string) (api.CanRegisterNodeResponse, error) {
	responseBytes, err := c.callAPI("node can-register", timezoneLocation)
	if err != nil {
		return api.CanRegisterNodeResponse{}, fmt.Errorf("Could not get can register node status: %w", err)
	}
	var response api.CanRegisterNodeResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanRegisterNodeResponse{}, fmt.Errorf("Could not decode can register node response: %w", err)
	}
	if response.Error != "" {
		return api.CanRegisterNodeResponse{}, fmt.Errorf("Could not get can register node status: %s", response.Error)
	}
	return response, nil
}

// Register the node
func (c *Client) RegisterNode(timezoneLocation string) (api.RegisterNodeResponse, error) {
	responseBytes, err := c.callAPI("node register", timezoneLocation)
	if err != nil {
		return api.RegisterNodeResponse{}, fmt.Errorf("Could not register node: %w", err)
	}
	var response api.RegisterNodeResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.RegisterNodeResponse{}, fmt.Errorf("Could not decode register node response: %w", err)
	}
	if response.Error != "" {
		return api.RegisterNodeResponse{}, fmt.Errorf("Could not register node: %s", response.Error)
	}
	return response, nil
}

// Checks if the node's withdrawal address can be set
func (c *Client) CanSetNodeWithdrawalAddress(withdrawalAddress common.Address, confirm bool) (api.CanSetNodeWithdrawalAddressResponse, error) {
	responseBytes, err := c.callAPI("node can-set-withdrawal-address", withdrawalAddress.Hex(), strconv.FormatBool(confirm))
	if err != nil {
		return api.CanSetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not get can set node withdrawal address: %w", err)
	}
	var response api.CanSetNodeWithdrawalAddressResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanSetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not decode can set node withdrawal address response: %w", err)
	}
	if response.Error != "" {
		return api.CanSetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not get can set node withdrawal address: %s", response.Error)
	}
	return response, nil
}

// Set the node's withdrawal address
func (c *Client) SetNodeWithdrawalAddress(withdrawalAddress common.Address, confirm bool) (api.SetNodeWithdrawalAddressResponse, error) {
	responseBytes, err := c.callAPI("node set-withdrawal-address", withdrawalAddress.Hex(), strconv.FormatBool(confirm))
	if err != nil {
		return api.SetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not set node withdrawal address: %w", err)
	}
	var response api.SetNodeWithdrawalAddressResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.SetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not decode set node withdrawal address response: %w", err)
	}
	if response.Error != "" {
		return api.SetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not set node withdrawal address: %s", response.Error)
	}
	return response, nil
}

// Checks if the node's withdrawal address can be confirmed
func (c *Client) CanConfirmNodeWithdrawalAddress() (api.CanSetNodeWithdrawalAddressResponse, error) {
	responseBytes, err := c.callAPI("node can-confirm-withdrawal-address")
	if err != nil {
		return api.CanSetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not get can confirm node withdrawal address: %w", err)
	}
	var response api.CanSetNodeWithdrawalAddressResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanSetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not decode can confirm node withdrawal address response: %w", err)
	}
	if response.Error != "" {
		return api.CanSetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not get can confirm node withdrawal address: %s", response.Error)
	}
	return response, nil
}

// Confirm the node's withdrawal address
func (c *Client) ConfirmNodeWithdrawalAddress() (api.SetNodeWithdrawalAddressResponse, error) {
	responseBytes, err := c.callAPI("node confirm-withdrawal-address")
	if err != nil {
		return api.SetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not confirm node withdrawal address: %w", err)
	}
	var response api.SetNodeWithdrawalAddressResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.SetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not decode confirm node withdrawal address response: %w", err)
	}
	if response.Error != "" {
		return api.SetNodeWithdrawalAddressResponse{}, fmt.Errorf("Could not confirm node withdrawal address: %s", response.Error)
	}
	return response, nil
}

// Checks if the node's timezone location can be set
func (c *Client) CanSetNodeTimezone(timezoneLocation string) (api.CanSetNodeTimezoneResponse, error) {
	responseBytes, err := c.callAPI("node can-set-timezone", timezoneLocation)
	if err != nil {
		return api.CanSetNodeTimezoneResponse{}, fmt.Errorf("Could not get can set node timezone: %w", err)
	}
	var response api.CanSetNodeTimezoneResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanSetNodeTimezoneResponse{}, fmt.Errorf("Could not decode can set node timezone response: %w", err)
	}
	if response.Error != "" {
		return api.CanSetNodeTimezoneResponse{}, fmt.Errorf("Could not get can set node timezone: %s", response.Error)
	}
	return response, nil
}

// Set the node's timezone location
func (c *Client) SetNodeTimezone(timezoneLocation string) (api.SetNodeTimezoneResponse, error) {
	responseBytes, err := c.callAPI("node set-timezone", timezoneLocation)
	if err != nil {
		return api.SetNodeTimezoneResponse{}, fmt.Errorf("Could not set node timezone: %w", err)
	}
	var response api.SetNodeTimezoneResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.SetNodeTimezoneResponse{}, fmt.Errorf("Could not decode set node timezone response: %w", err)
	}
	if response.Error != "" {
		return api.SetNodeTimezoneResponse{}, fmt.Errorf("Could not set node timezone: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can swap GGP tokens
func (c *Client) CanNodeSwapGgp(amountWei *big.Int) (api.CanNodeSwapGgpResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node can-swap-ggp %s", amountWei.String()))
	if err != nil {
		return api.CanNodeSwapGgpResponse{}, fmt.Errorf("Could not get can node swap GGP status: %w", err)
	}
	var response api.CanNodeSwapGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanNodeSwapGgpResponse{}, fmt.Errorf("Could not decode can node swap GGP response: %w", err)
	}
	if response.Error != "" {
		return api.CanNodeSwapGgpResponse{}, fmt.Errorf("Could not get can node swap GGP status: %s", response.Error)
	}
	return response, nil
}

// Get the gas estimate for approving legacy GGP interaction
func (c *Client) NodeSwapGgpApprovalGas(amountWei *big.Int) (api.NodeSwapGgpApproveGasResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node get-swap-ggp-approval-gas %s", amountWei.String()))
	if err != nil {
		return api.NodeSwapGgpApproveGasResponse{}, fmt.Errorf("Could not get old GGP approval gas: %w", err)
	}
	var response api.NodeSwapGgpApproveGasResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeSwapGgpApproveGasResponse{}, fmt.Errorf("Could not decode node swap GGP approve gas response: %w", err)
	}
	if response.Error != "" {
		return api.NodeSwapGgpApproveGasResponse{}, fmt.Errorf("Could not get old GGP approval gas: %s", response.Error)
	}
	return response, nil
}

// Approves old GGP for a token swap
func (c *Client) NodeSwapGgpApprove(amountWei *big.Int) (api.NodeSwapGgpApproveResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node swap-ggp-approve-ggp %s", amountWei.String()))
	if err != nil {
		return api.NodeSwapGgpApproveResponse{}, fmt.Errorf("Could not approve old GGP: %w", err)
	}
	var response api.NodeSwapGgpApproveResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeSwapGgpApproveResponse{}, fmt.Errorf("Could not decode node swap GGP approve response: %w", err)
	}
	if response.Error != "" {
		return api.NodeSwapGgpApproveResponse{}, fmt.Errorf("Could not approve old GGP tokens for swapping: %s", response.Error)
	}
	return response, nil
}

// Swap node's old GGP tokens for new GGP tokens, waiting for the approval to be mined first
func (c *Client) NodeWaitAndSwapGgp(amountWei *big.Int, approvalTxHash common.Hash) (api.NodeSwapGgpSwapResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node wait-and-swap-ggp %s %s", amountWei.String(), approvalTxHash.String()))
	if err != nil {
		return api.NodeSwapGgpSwapResponse{}, fmt.Errorf("Could not swap node's GGP tokens: %w", err)
	}
	var response api.NodeSwapGgpSwapResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeSwapGgpSwapResponse{}, fmt.Errorf("Could not decode node swap GGP tokens response: %w", err)
	}
	if response.Error != "" {
		return api.NodeSwapGgpSwapResponse{}, fmt.Errorf("Could not swap node's GGP tokens: %s", response.Error)
	}
	return response, nil
}

// Swap node's old GGP tokens for new GGP tokens
func (c *Client) NodeSwapGgp(amountWei *big.Int) (api.NodeSwapGgpSwapResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node swap-ggp %s", amountWei.String()))
	if err != nil {
		return api.NodeSwapGgpSwapResponse{}, fmt.Errorf("Could not swap node's GGP tokens: %w", err)
	}
	var response api.NodeSwapGgpSwapResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeSwapGgpSwapResponse{}, fmt.Errorf("Could not decode node swap GGP tokens response: %w", err)
	}
	if response.Error != "" {
		return api.NodeSwapGgpSwapResponse{}, fmt.Errorf("Could not swap node's GGP tokens: %s", response.Error)
	}
	return response, nil
}

// Get a node's legacy GGP allowance for swapping on the new GGP contract
func (c *Client) GetNodeSwapGgpAllowance() (api.NodeSwapGgpAllowanceResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node swap-ggp-allowance"))
	if err != nil {
		return api.NodeSwapGgpAllowanceResponse{}, fmt.Errorf("Could not get node swap GGP allowance: %w", err)
	}
	var response api.NodeSwapGgpAllowanceResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeSwapGgpAllowanceResponse{}, fmt.Errorf("Could not decode node swap GGP allowance response: %w", err)
	}
	if response.Error != "" {
		return api.NodeSwapGgpAllowanceResponse{}, fmt.Errorf("Could not get node swap GGP allowance: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can stake GGP
func (c *Client) CanNodeStakeGgp(amountWei *big.Int) (api.CanNodeStakeGgpResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node can-stake-ggp %s", amountWei.String()))
	if err != nil {
		return api.CanNodeStakeGgpResponse{}, fmt.Errorf("Could not get can node stake GGP status: %w", err)
	}
	var response api.CanNodeStakeGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanNodeStakeGgpResponse{}, fmt.Errorf("Could not decode can node stake GGP response: %w", err)
	}
	if response.Error != "" {
		return api.CanNodeStakeGgpResponse{}, fmt.Errorf("Could not get can node stake GGP status: %s", response.Error)
	}
	return response, nil
}

// Get the gas estimate for approving new GGP interaction
func (c *Client) NodeStakeGgpApprovalGas(amountWei *big.Int) (api.NodeStakeGgpApproveGasResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node get-stake-ggp-approval-gas %s", amountWei.String()))
	if err != nil {
		return api.NodeStakeGgpApproveGasResponse{}, fmt.Errorf("Could not get new GGP approval gas: %w", err)
	}
	var response api.NodeStakeGgpApproveGasResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeStakeGgpApproveGasResponse{}, fmt.Errorf("Could not decode node stake GGP approve gas response: %w", err)
	}
	if response.Error != "" {
		return api.NodeStakeGgpApproveGasResponse{}, fmt.Errorf("Could not get new GGP approval gas: %s", response.Error)
	}
	return response, nil
}

// Approve GGP for staking against the node
func (c *Client) NodeStakeGgpApprove(amountWei *big.Int) (api.NodeStakeGgpApproveResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node stake-ggp-approve-ggp %s", amountWei.String()))
	if err != nil {
		return api.NodeStakeGgpApproveResponse{}, fmt.Errorf("Could not approve GGP for staking: %w", err)
	}
	var response api.NodeStakeGgpApproveResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeStakeGgpApproveResponse{}, fmt.Errorf("Could not decode stake node GGP approve response: %w", err)
	}
	if response.Error != "" {
		return api.NodeStakeGgpApproveResponse{}, fmt.Errorf("Could not approve GGP for staking: %s", response.Error)
	}
	return response, nil
}

// Stake GGP against the node waiting for approvalTxHash to be mined first
func (c *Client) NodeWaitAndStakeGgp(amountWei *big.Int, approvalTxHash common.Hash) (api.NodeStakeGgpStakeResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node wait-and-stake-ggp %s %s", amountWei.String(), approvalTxHash.String()))
	if err != nil {
		return api.NodeStakeGgpStakeResponse{}, fmt.Errorf("Could not stake node GGP: %w", err)
	}
	var response api.NodeStakeGgpStakeResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeStakeGgpStakeResponse{}, fmt.Errorf("Could not decode stake node GGP response: %w", err)
	}
	if response.Error != "" {
		return api.NodeStakeGgpStakeResponse{}, fmt.Errorf("Could not stake node GGP: %s", response.Error)
	}
	return response, nil
}

// Stake GGP against the node
func (c *Client) NodeStakeGgp(amountWei *big.Int) (api.NodeStakeGgpStakeResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node stake-ggp %s", amountWei.String()))
	if err != nil {
		return api.NodeStakeGgpStakeResponse{}, fmt.Errorf("Could not stake node GGP: %w", err)
	}
	var response api.NodeStakeGgpStakeResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeStakeGgpStakeResponse{}, fmt.Errorf("Could not decode stake node GGP response: %w", err)
	}
	if response.Error != "" {
		return api.NodeStakeGgpStakeResponse{}, fmt.Errorf("Could not stake node GGP: %s", response.Error)
	}
	return response, nil
}

// Get a node's GGP allowance for the staking contract
func (c *Client) GetNodeStakeGgpAllowance() (api.NodeStakeGgpAllowanceResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node stake-ggp-allowance"))
	if err != nil {
		return api.NodeStakeGgpAllowanceResponse{}, fmt.Errorf("Could not get node stake GGP allowance: %w", err)
	}
	var response api.NodeStakeGgpAllowanceResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeStakeGgpAllowanceResponse{}, fmt.Errorf("Could not decode node stake GGP allowance response: %w", err)
	}
	if response.Error != "" {
		return api.NodeStakeGgpAllowanceResponse{}, fmt.Errorf("Could not get node stake GGP allowance: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can withdraw GGP
func (c *Client) CanNodeWithdrawGgp(amountWei *big.Int) (api.CanNodeWithdrawGgpResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node can-withdraw-ggp %s", amountWei.String()))
	if err != nil {
		return api.CanNodeWithdrawGgpResponse{}, fmt.Errorf("Could not get can node withdraw GGP status: %w", err)
	}
	var response api.CanNodeWithdrawGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanNodeWithdrawGgpResponse{}, fmt.Errorf("Could not decode can node withdraw GGP response: %w", err)
	}
	if response.Error != "" {
		return api.CanNodeWithdrawGgpResponse{}, fmt.Errorf("Could not get can node withdraw GGP status: %s", response.Error)
	}
	return response, nil
}

// Withdraw GGP staked against the node
func (c *Client) NodeWithdrawGgp(amountWei *big.Int) (api.NodeWithdrawGgpResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node withdraw-ggp %s", amountWei.String()))
	if err != nil {
		return api.NodeWithdrawGgpResponse{}, fmt.Errorf("Could not withdraw node GGP: %w", err)
	}
	var response api.NodeWithdrawGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeWithdrawGgpResponse{}, fmt.Errorf("Could not decode withdraw node GGP response: %w", err)
	}
	if response.Error != "" {
		return api.NodeWithdrawGgpResponse{}, fmt.Errorf("Could not withdraw node GGP: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can make a deposit
func (c *Client) CanNodeDeposit(amountWei *big.Int, minFee float64, salt *big.Int) (api.CanNodeDepositResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node can-deposit %s %f %s", amountWei.String(), minFee, salt.String()))
	if err != nil {
		return api.CanNodeDepositResponse{}, fmt.Errorf("Could not get can node deposit status: %w", err)
	}
	var response api.CanNodeDepositResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanNodeDepositResponse{}, fmt.Errorf("Could not decode can node deposit response: %w", err)
	}
	if response.Error != "" {
		return api.CanNodeDepositResponse{}, fmt.Errorf("Could not get can node deposit status: %s", response.Error)
	}
	return response, nil
}

// Make a node deposit
func (c *Client) NodeDeposit(amountWei *big.Int, minFee float64, salt *big.Int) (api.NodeDepositResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node deposit %s %f %s", amountWei.String(), minFee, salt.String()))
	if err != nil {
		return api.NodeDepositResponse{}, fmt.Errorf("Could not make node deposit: %w", err)
	}
	var response api.NodeDepositResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeDepositResponse{}, fmt.Errorf("Could not decode node deposit response: %w", err)
	}
	if response.Error != "" {
		return api.NodeDepositResponse{}, fmt.Errorf("Could not make node deposit: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can send tokens
func (c *Client) CanNodeSend(amountWei *big.Int, token string) (api.CanNodeSendResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node can-send %s %s", amountWei.String(), token))
	if err != nil {
		return api.CanNodeSendResponse{}, fmt.Errorf("Could not get can node send status: %w", err)
	}
	var response api.CanNodeSendResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanNodeSendResponse{}, fmt.Errorf("Could not decode can node send response: %w", err)
	}
	if response.Error != "" {
		return api.CanNodeSendResponse{}, fmt.Errorf("Could not get can node send status: %s", response.Error)
	}
	return response, nil
}

// Send tokens from the node to an address
func (c *Client) NodeSend(amountWei *big.Int, token string, toAddress common.Address) (api.NodeSendResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node send %s %s %s", amountWei.String(), token, toAddress.Hex()))
	if err != nil {
		return api.NodeSendResponse{}, fmt.Errorf("Could not send tokens from node: %w", err)
	}
	var response api.NodeSendResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeSendResponse{}, fmt.Errorf("Could not decode node send response: %w", err)
	}
	if response.Error != "" {
		return api.NodeSendResponse{}, fmt.Errorf("Could not send tokens from node: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can burn tokens
func (c *Client) CanNodeBurn(amountWei *big.Int, token string) (api.CanNodeBurnResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node can-burn %s %s", amountWei.String(), token))
	if err != nil {
		return api.CanNodeBurnResponse{}, fmt.Errorf("Could not get can node burn status: %w", err)
	}
	var response api.CanNodeBurnResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanNodeBurnResponse{}, fmt.Errorf("Could not decode can node burn response: %w", err)
	}
	if response.Error != "" {
		return api.CanNodeBurnResponse{}, fmt.Errorf("Could not get can node burn status: %s", response.Error)
	}
	return response, nil
}

// Burn tokens owned by the node for ETH
func (c *Client) NodeBurn(amountWei *big.Int, token string) (api.NodeBurnResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("node burn %s %s", amountWei.String(), token))
	if err != nil {
		return api.NodeBurnResponse{}, fmt.Errorf("Could not burn tokens owned by node: %w", err)
	}
	var response api.NodeBurnResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeBurnResponse{}, fmt.Errorf("Could not decode node burn response: %w", err)
	}
	if response.Error != "" {
		return api.NodeBurnResponse{}, fmt.Errorf("Could not burn tokens owned by node: %s", response.Error)
	}
	return response, nil
}

// Get node sync progress
func (c *Client) NodeSync() (api.NodeSyncProgressResponse, error) {
	responseBytes, err := c.callAPI("node sync")
	if err != nil {
		return api.NodeSyncProgressResponse{}, fmt.Errorf("Could not get node sync: %w", err)
	}
	var response api.NodeSyncProgressResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeSyncProgressResponse{}, fmt.Errorf("Could not decode node sync response: %w", err)
	}
	if response.Error != "" {
		return api.NodeSyncProgressResponse{}, fmt.Errorf("Could not get node sync: %s", response.Error)
	}
	return response, nil
}

// Check whether the node has GGP rewards available to claim
func (c *Client) CanNodeClaimGgp() (api.CanNodeClaimGgpResponse, error) {
	responseBytes, err := c.callAPI("node can-claim-ggp-rewards")
	if err != nil {
		return api.CanNodeClaimGgpResponse{}, fmt.Errorf("Could not get can node claim ggp rewards status: %w", err)
	}
	var response api.CanNodeClaimGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanNodeClaimGgpResponse{}, fmt.Errorf("Could not decode can node claim ggp rewards response: %w", err)
	}
	if response.Error != "" {
		return api.CanNodeClaimGgpResponse{}, fmt.Errorf("Could not get can node claim ggp rewards status: %s", response.Error)
	}
	return response, nil
}

// Claim available GGP rewards
func (c *Client) NodeClaimGgp() (api.NodeClaimGgpResponse, error) {
	responseBytes, err := c.callAPI("node claim-ggp-rewards")
	if err != nil {
		return api.NodeClaimGgpResponse{}, fmt.Errorf("Could not claim ggp rewards: %w", err)
	}
	var response api.NodeClaimGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeClaimGgpResponse{}, fmt.Errorf("Could not decode node claim ggp rewards response: %w", err)
	}
	if response.Error != "" {
		return api.NodeClaimGgpResponse{}, fmt.Errorf("Could not claim ggp rewards: %s", response.Error)
	}
	return response, nil
}

// Get node GGP rewards status
func (c *Client) NodeRewards() (api.NodeRewardsResponse, error) {
	responseBytes, err := c.callAPI("node rewards")
	if err != nil {
		return api.NodeRewardsResponse{}, fmt.Errorf("Could not get node rewards: %w", err)
	}
	var response api.NodeRewardsResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.NodeRewardsResponse{}, fmt.Errorf("Could not decode node rewards response: %w", err)
	}
	if response.Error != "" {
		return api.NodeRewardsResponse{}, fmt.Errorf("Could not get node rewards: %s", response.Error)
	}
	return response, nil
}

// Get the deposit contract info for Rocket Pool and the Beacon Client
func (c *Client) DepositContractInfo() (api.DepositContractInfoResponse, error) {
	responseBytes, err := c.callAPI("node deposit-contract-info")
	if err != nil {
		return api.DepositContractInfoResponse{}, fmt.Errorf("Could not get deposit contract info: %w", err)
	}
	var response api.DepositContractInfoResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.DepositContractInfoResponse{}, fmt.Errorf("Could not decode deposit contract info response: %w", err)
	}
	if response.Error != "" {
		return api.DepositContractInfoResponse{}, fmt.Errorf("Could not get deposit contract info: %s", response.Error)
	}
	return response, nil
}
