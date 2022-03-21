package auction

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/auction"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
)

func canRecoverGgpFromLot(c *cli.Context, lotIndex uint64) (*api.CanRecoverGGPFromLotResponse, error) {

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
	response := api.CanRecoverGGPFromLotResponse{}

	// Sync
	var wg errgroup.Group

	// Check if lot exists
	wg.Go(func() error {
		lotExists, err := auction.GetLotExists(rp, lotIndex, nil)
		if err == nil {
			response.DoesNotExist = !lotExists
		}
		return err
	})

	// Check if lot bidding has ended
	wg.Go(func() error {
		biddingEnded, err := getLotBiddingEnded(rp, lotIndex)
		if err == nil {
			response.BiddingNotEnded = !biddingEnded
		}
		return err
	})

	// Check if lot contains unclaimed GGP
	wg.Go(func() error {
		remainingGgp, err := auction.GetLotRemainingGGPAmount(rp, lotIndex, nil)
		if err == nil {
			response.NoUnclaimedGGP = (remainingGgp.Cmp(big.NewInt(0)) == 0)
		}
		return err
	})

	// Check if unclaimed GGP has already been recovered
	wg.Go(func() error {
		ggpRecovered, err := auction.GetLotGGPRecovered(rp, lotIndex, nil)
		if err == nil {
			response.GGPAlreadyRecovered = ggpRecovered
		}
		return err
	})

	// Get gas estimate
	wg.Go(func() error {
		opts, err := w.GetNodeAccountTransactor()
		if err != nil {
			return err
		}
		gasInfo, err := auction.EstimateRecoverUnclaimedGGPGas(rp, lotIndex, opts)
		if err == nil {
			response.GasInfo = gasInfo
		}
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Update & return response
	response.CanRecover = !(response.DoesNotExist || response.BiddingNotEnded || response.NoUnclaimedGGP || response.GGPAlreadyRecovered)
	return &response, nil

}

func recoverGgpFromLot(c *cli.Context, lotIndex uint64) (*api.RecoverGGPFromLotResponse, error) {

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
	response := api.RecoverGGPFromLotResponse{}

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

	// Recover unclaimed GGP from lot
	hash, err := auction.RecoverUnclaimedGGP(rp, lotIndex, opts)
	if err != nil {
		return nil, err
	}
	response.TxHash = hash

	// Return response
	return &response, nil

}
