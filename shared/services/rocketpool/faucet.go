package rocketpool

import (
	"encoding/json"
	"fmt"

	"github.com/rocket-pool/smartnode/shared/types/api"
)

// Get faucet status
func (c *Client) FaucetStatus() (api.FaucetStatusResponse, error) {
	responseBytes, err := c.callAPI("faucet status")
	if err != nil {
		return api.FaucetStatusResponse{}, fmt.Errorf("Could not get faucet status: %w", err)
	}
	var response api.FaucetStatusResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.FaucetStatusResponse{}, fmt.Errorf("Could not decode faucet status response: %w", err)
	}
	if response.Error != "" {
		return api.FaucetStatusResponse{}, fmt.Errorf("Could not get faucet status: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can withdraw GGP from the faucet
func (c *Client) CanFaucetWithdrawGgp() (api.CanFaucetWithdrawGgpResponse, error) {
	responseBytes, err := c.callAPI("faucet can-withdraw-ggp")
	if err != nil {
		return api.CanFaucetWithdrawGgpResponse{}, fmt.Errorf("Could not get can withdraw GGP from faucet status: %w", err)
	}
	var response api.CanFaucetWithdrawGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanFaucetWithdrawGgpResponse{}, fmt.Errorf("Could not decode can withdraw GGP from faucet response: %w", err)
	}
	if response.Error != "" {
		return api.CanFaucetWithdrawGgpResponse{}, fmt.Errorf("Could not get can withdraw GGP from faucet status: %s", response.Error)
	}
	return response, nil
}

// Withdraw GGP from the faucet
func (c *Client) FaucetWithdrawGgp() (api.FaucetWithdrawGgpResponse, error) {
	responseBytes, err := c.callAPI("faucet withdraw-ggp")
	if err != nil {
		return api.FaucetWithdrawGgpResponse{}, fmt.Errorf("Could not withdraw GGP from faucet: %w", err)
	}
	var response api.FaucetWithdrawGgpResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.FaucetWithdrawGgpResponse{}, fmt.Errorf("Could not decode withdraw GGP from faucet response: %w", err)
	}
	if response.Error != "" {
		return api.FaucetWithdrawGgpResponse{}, fmt.Errorf("Could not withdraw GGP from faucet: %s", response.Error)
	}
	return response, nil
}
