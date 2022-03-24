package auction

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)

func getLots(c *cli.Context) error {

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Get lot details
	lots, err := rp.AuctionLots()
	if err != nil {
		return err
	}

	// Get lots by status
	openLots := []api.LotDetails{}
	clearedLots := []api.LotDetails{}
	claimableLots := []api.LotDetails{}
	biddableLots := []api.LotDetails{}
	recoverableLots := []api.LotDetails{}
	for _, lot := range lots.Lots {
		if lot.Details.Cleared {
			clearedLots = append(clearedLots, lot)
		} else {
			openLots = append(openLots, lot)
		}
		if lot.ClaimAvailable {
			claimableLots = append(claimableLots, lot)
		}
		if lot.BiddingAvailable {
			biddableLots = append(biddableLots, lot)
		}
		if lot.GGPRecoveryAvailable {
			recoverableLots = append(recoverableLots, lot)
		}
	}

	// Print lot details by status
	if len(lots.Lots) == 0 {
		fmt.Println("There are no lots for auction yet.")
	}
	for status := 0; status < 2; status++ {

		// Get status title format & lot list
		var statusFormat string
		var statusLots []api.LotDetails
		if status == 0 {
			statusFormat = "%d lot(s) open for bidding:\n"
			statusLots = openLots
		} else {
			statusFormat = "%d cleared lot(s):\n"
			statusLots = clearedLots
		}
		if len(statusLots) == 0 {
			continue
		}

		// Print
		fmt.Printf(statusFormat, len(statusLots))
		for _, lot := range statusLots {
			fmt.Printf("--------------------\n")
			fmt.Printf("\n")
			fmt.Printf("Lot ID:               %d\n", lot.Details.Index)
			fmt.Printf("Start block:          %d\n", lot.Details.StartBlock)
			fmt.Printf("End block:            %d\n", lot.Details.EndBlock)
			fmt.Printf("GGP starting price:   %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.StartPrice), 6))
			fmt.Printf("GGP reserve price:    %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.ReservePrice), 6))
			fmt.Printf("GGP current price:    %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.CurrentPrice), 6))
			fmt.Printf("Total GGP amount:     %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.TotalGGPAmount), 6))
			fmt.Printf("Claimed GGP amount:   %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.ClaimedGGPAmount), 6))
			fmt.Printf("Remaining GGP amount: %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.RemainingGGPAmount), 6))
			fmt.Printf("Total ETH bid:        %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.TotalBidAmount), 6))
			fmt.Printf("ETH bid by node:      %.6f\n", math.RoundDown(eth.WeiToEth(lot.Details.AddressBidAmount), 6))
			if lot.Details.Cleared {
				fmt.Printf("Cleared:              yes\n")
				if lot.Details.RemainingGGPAmount.Cmp(big.NewInt(0)) == 0 {
					fmt.Printf("Unclaimed GGP:        no\n")
				} else if lot.Details.GGPRecovered {
					fmt.Printf("Unclaimed GGP:        recovered\n")
				} else {
					fmt.Printf("Unclaimed GGP:        yes\n")
				}
			} else {
				fmt.Printf("Cleared:              no\n")
			}
			fmt.Printf("\n")
		}
		fmt.Println("")

	}

	// Print actionable lot details
	if len(claimableLots) > 0 {
		fmt.Printf("%d lot(s) you have bid on have GGP available to claim:\n", len(claimableLots))
		for _, lot := range claimableLots {
			fmt.Printf("- lot %d (%.6f ETH bid @ %.6f ETH per GGP)\n", lot.Details.Index, math.RoundDown(eth.WeiToEth(lot.Details.AddressBidAmount), 6), math.RoundDown(eth.WeiToEth(lot.Details.CurrentPrice), 6))
		}
		fmt.Println("")
	}
	if len(biddableLots) > 0 {
		fmt.Printf("%d lot(s) are open for bidding:\n", len(biddableLots))
		for _, lot := range biddableLots {
			fmt.Printf("- lot %d (%.6f GGP available @ %.6f ETH per GGP)\n", lot.Details.Index, math.RoundDown(eth.WeiToEth(lot.Details.RemainingGGPAmount), 6), math.RoundDown(eth.WeiToEth(lot.Details.CurrentPrice), 6))
		}
		fmt.Println("")
	}
	if len(recoverableLots) > 0 {
		fmt.Printf("%d lot(s) have unclaimed GGP ready to recover:\n", len(recoverableLots))
		for _, lot := range recoverableLots {
			fmt.Printf("- lot %d (%.6f GGP unclaimed)\n", lot.Details.Index, math.RoundDown(eth.WeiToEth(lot.Details.RemainingGGPAmount), 6))
		}
		fmt.Println("")
	}

	// Return
	return nil

}
