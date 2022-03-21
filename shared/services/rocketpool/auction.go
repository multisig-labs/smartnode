package rocketpool

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/rocket-pool/smartnode/shared/types/api"
)

// Get GGP auction status
func (c *Client) AuctionStatus() (api.AuctionStatusResponse, error) {
	responseBytes, err := c.callAPI("auction status")
	if err != nil {
		return api.AuctionStatusResponse{}, fmt.Errorf("Could not get auction status: %w", err)
	}
	var response api.AuctionStatusResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.AuctionStatusResponse{}, fmt.Errorf("Could not decode auction stats response: %w", err)
	}
	if response.Error != "" {
		return api.AuctionStatusResponse{}, fmt.Errorf("Could not get auction status: %s", response.Error)
	}
	if response.TotalGGPBalance == nil {
		response.TotalGGPBalance = big.NewInt(0)
	}
	if response.AllottedGGPBalance == nil {
		response.AllottedGGPBalance = big.NewInt(0)
	}
	if response.RemainingGGPBalance == nil {
		response.RemainingGGPBalance = big.NewInt(0)
	}
	return response, nil
}

// Get GGP lots for auction
func (c *Client) AuctionLots() (api.AuctionLotsResponse, error) {
	responseBytes, err := c.callAPI("auction lots")
	if err != nil {
		return api.AuctionLotsResponse{}, fmt.Errorf("Could not get auction lots: %w", err)
	}
	var response api.AuctionLotsResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.AuctionLotsResponse{}, fmt.Errorf("Could not decode auction lots response: %w", err)
	}
	if response.Error != "" {
		return api.AuctionLotsResponse{}, fmt.Errorf("Could not get auction lots: %s", response.Error)
	}
	for i := 0; i < len(response.Lots); i++ {
		details := &response.Lots[i].Details
		if details.StartPrice == nil {
			details.StartPrice = big.NewInt(0)
		}
		if details.ReservePrice == nil {
			details.ReservePrice = big.NewInt(0)
		}
		if details.PriceAtCurrentBlock == nil {
			details.PriceAtCurrentBlock = big.NewInt(0)
		}
		if details.PriceByTotalBids == nil {
			details.PriceByTotalBids = big.NewInt(0)
		}
		if details.CurrentPrice == nil {
			details.CurrentPrice = big.NewInt(0)
		}
		if details.TotalGGPAmount == nil {
			details.TotalGGPAmount = big.NewInt(0)
		}
		if details.ClaimedGGPAmount == nil {
			details.ClaimedGGPAmount = big.NewInt(0)
		}
		if details.RemainingGGPAmount == nil {
			details.RemainingGGPAmount = big.NewInt(0)
		}
		if details.TotalBidAmount == nil {
			details.TotalBidAmount = big.NewInt(0)
		}
		if details.AddressBidAmount == nil {
			details.AddressBidAmount = big.NewInt(0)
		}
	}
	return response, nil
}

// Check whether the node can create a new lot
func (c *Client) CanCreateLot() (api.CanCreateLotResponse, error) {
	responseBytes, err := c.callAPI("auction can-create-lot")
	if err != nil {
		return api.CanCreateLotResponse{}, fmt.Errorf("Could not get can create lot status: %w", err)
	}
	var response api.CanCreateLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanCreateLotResponse{}, fmt.Errorf("Could not decode can create lot response: %w", err)
	}
	if response.Error != "" {
		return api.CanCreateLotResponse{}, fmt.Errorf("Could not get can create lot status: %s", response.Error)
	}
	return response, nil
}

// Create a new lot
func (c *Client) CreateLot() (api.CreateLotResponse, error) {
	responseBytes, err := c.callAPI("auction create-lot")
	if err != nil {
		return api.CreateLotResponse{}, fmt.Errorf("Could not create lot: %w", err)
	}
	var response api.CreateLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CreateLotResponse{}, fmt.Errorf("Could not decode create lot response: %w", err)
	}
	if response.Error != "" {
		return api.CreateLotResponse{}, fmt.Errorf("Could not create lot: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can bid on a lot
func (c *Client) CanBidOnLot(lotIndex uint64, amountWei *big.Int) (api.CanBidOnLotResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("auction can-bid-lot %d %s", lotIndex, amountWei.String()))
	if err != nil {
		return api.CanBidOnLotResponse{}, fmt.Errorf("Could not get can bid on lot status: %w", err)
	}
	var response api.CanBidOnLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanBidOnLotResponse{}, fmt.Errorf("Could not decode can bid on lot response: %w", err)
	}
	if response.Error != "" {
		return api.CanBidOnLotResponse{}, fmt.Errorf("Could not get can bid on lot status: %s", response.Error)
	}
	return response, nil
}

// Bid on a lot
func (c *Client) BidOnLot(lotIndex uint64, amountWei *big.Int) (api.BidOnLotResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("auction bid-lot %d %s", lotIndex, amountWei.String()))
	if err != nil {
		return api.BidOnLotResponse{}, fmt.Errorf("Could not bid on lot: %w", err)
	}
	var response api.BidOnLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.BidOnLotResponse{}, fmt.Errorf("Could not decode bid on lot response: %w", err)
	}
	if response.Error != "" {
		return api.BidOnLotResponse{}, fmt.Errorf("Could not bid on lot: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can claim GGP from a lot
func (c *Client) CanClaimFromLot(lotIndex uint64) (api.CanClaimFromLotResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("auction can-claim-lot %d", lotIndex))
	if err != nil {
		return api.CanClaimFromLotResponse{}, fmt.Errorf("Could not get can claim GGP from lot status: %w", err)
	}
	var response api.CanClaimFromLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanClaimFromLotResponse{}, fmt.Errorf("Could not decode can claim GGP from lot response: %w", err)
	}
	if response.Error != "" {
		return api.CanClaimFromLotResponse{}, fmt.Errorf("Could not get can claim GGP from lot status: %s", response.Error)
	}
	return response, nil
}

// Claim GGP from a lot
func (c *Client) ClaimFromLot(lotIndex uint64) (api.ClaimFromLotResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("auction claim-lot %d", lotIndex))
	if err != nil {
		return api.ClaimFromLotResponse{}, fmt.Errorf("Could not claim GGP from lot: %w", err)
	}
	var response api.ClaimFromLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.ClaimFromLotResponse{}, fmt.Errorf("Could not decode claim GGP from lot response: %w", err)
	}
	if response.Error != "" {
		return api.ClaimFromLotResponse{}, fmt.Errorf("Could not claim GGP from lot: %s", response.Error)
	}
	return response, nil
}

// Check whether the node can recover unclaimed GGP from a lot
func (c *Client) CanRecoverUnclaimedGGPFromLot(lotIndex uint64) (api.CanRecoverGGPFromLotResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("auction can-recover-lot %d", lotIndex))
	if err != nil {
		return api.CanRecoverGGPFromLotResponse{}, fmt.Errorf("Could not get can recover unclaimed GGP from lot status: %w", err)
	}
	var response api.CanRecoverGGPFromLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.CanRecoverGGPFromLotResponse{}, fmt.Errorf("Could not decode can recover unclaimed GGP from lot response: %w", err)
	}
	if response.Error != "" {
		return api.CanRecoverGGPFromLotResponse{}, fmt.Errorf("Could not get can recover unclaimed GGP from lot status: %s", response.Error)
	}
	return response, nil
}

// Recover unclaimed GGP from a lot (returning it to the auction contract)
func (c *Client) RecoverUnclaimedGGPFromLot(lotIndex uint64) (api.RecoverGGPFromLotResponse, error) {
	responseBytes, err := c.callAPI(fmt.Sprintf("auction recover-lot %d", lotIndex))
	if err != nil {
		return api.RecoverGGPFromLotResponse{}, fmt.Errorf("Could not recover unclaimed GGP from lot: %w", err)
	}
	var response api.RecoverGGPFromLotResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return api.RecoverGGPFromLotResponse{}, fmt.Errorf("Could not decode recover unclaimed GGP from lot response: %w", err)
	}
	if response.Error != "" {
		return api.RecoverGGPFromLotResponse{}, fmt.Errorf("Could not recover unclaimed GGP from lot: %s", response.Error)
	}
	return response, nil
}
