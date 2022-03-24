package auction

import (
	"github.com/rocket-pool/rocketpool-go/auction"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
)

func getStatus(c *cli.Context) (*api.AuctionStatusResponse, error) {

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
	response := api.AuctionStatusResponse{}

	// Sync
	var wg errgroup.Group

	// Get auction contract GGP balances
	wg.Go(func() error {
		totalGgpBalance, err := auction.GetTotalGGPBalance(rp, nil)
		if err == nil {
			response.TotalGGPBalance = totalGgpBalance
		}
		return err
	})
	wg.Go(func() error {
		allottedGgpBalance, err := auction.GetAllottedGGPBalance(rp, nil)
		if err == nil {
			response.AllottedGGPBalance = allottedGgpBalance
		}
		return err
	})
	wg.Go(func() error {
		remainingGgpBalance, err := auction.GetRemainingGGPBalance(rp, nil)
		if err == nil {
			response.RemainingGGPBalance = remainingGgpBalance
		}
		return err
	})

	// Check if lot can be created
	wg.Go(func() error {
		sufficientRemainingGgpForLot, err := getSufficientRemainingGGPForLot(rp)
		if err == nil {
			response.CanCreateLot = sufficientRemainingGgpForLot
		}
		return err
	})

	// Get lot counts
	wg.Go(func() error {
		nodeAccount, err := w.GetNodeAccount()
		if err != nil {
			return err
		}
		lotCountDetails, err := getAllLotCountDetails(rp, nodeAccount.Address)
		if err == nil {
			for _, details := range lotCountDetails {
				if details.AddressHasBid && details.Cleared {
					response.LotCounts.ClaimAvailable++
				}
				if !details.Cleared && details.HasRemainingGgp {
					response.LotCounts.BiddingAvailable++
				}
				if details.Cleared && details.HasRemainingGgp && !details.GgpRecovered {
					response.LotCounts.GGPRecoveryAvailable++
				}
			}
		}
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Return response
	return &response, nil

}
