package auction

import (
	"fmt"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)

func getStatus(c *cli.Context) error {

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Get auction status
	status, err := rp.AuctionStatus()
	if err != nil {
		return err
	}

	// Print & return
	fmt.Printf(
		"A total of %.6f GGP is up for auction, with %.6f GGP currently allotted and %.6f GGP remaining.\n",
		math.RoundDown(eth.WeiToEth(status.TotalGGPBalance), 6),
		math.RoundDown(eth.WeiToEth(status.AllottedGGPBalance), 6),
		math.RoundDown(eth.WeiToEth(status.RemainingGGPBalance), 6))
	if status.LotCounts.ClaimAvailable > 0 {
		fmt.Printf("%d lot(s) you have bid on have GGP available to claim!\n", status.LotCounts.ClaimAvailable)
	}
	if status.LotCounts.BiddingAvailable > 0 {
		fmt.Printf("%d lot(s) are open for bidding!\n", status.LotCounts.BiddingAvailable)
	}
	if status.LotCounts.GGPRecoveryAvailable > 0 {
		fmt.Printf("%d cleared lot(s) have unclaimed GGP ready to recover!\n", status.LotCounts.GGPRecoveryAvailable)
	}
	if status.CanCreateLot {
		fmt.Println("A new lot can be created with remaining GGP in the auction contract.")
	}
	return nil

}
