package auction

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/auction"
	"github.com/rocket-pool/rocketpool-go/network"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/types/api"
)

// Settings
const LotCountDetailsBatchSize = 10
const LotDetailsBatchSize = 10

// Lot count details
type lotCountDetails struct {
	AddressHasBid   bool
	Cleared         bool
	HasRemainingGgp bool
	GgpRecovered    bool
}

// Check if bidding has ended for a lot
func getLotBiddingEnded(rp *rocketpool.RocketPool, lotIndex uint64) (bool, error) {

	// Data
	var wg errgroup.Group
	var currentBlock uint64
	var lotEndBlock uint64

	// Get current block
	wg.Go(func() error {
		header, err := rp.Client.HeaderByNumber(context.Background(), nil)
		if err == nil {
			currentBlock = header.Number.Uint64()
		}
		return err
	})

	// Get lot end block
	wg.Go(func() error {
		var err error
		lotEndBlock, err = auction.GetLotEndBlock(rp, lotIndex, nil)
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return false, err
	}

	// Return
	return (currentBlock >= lotEndBlock), nil

}

// Check whether sufficient remaining GGP is available to create a lot
func getSufficientRemainingGGPForLot(rp *rocketpool.RocketPool) (bool, error) {

	// Data
	var wg errgroup.Group
	var remainingGgpBalance *big.Int
	var lotMinimumEthValue *big.Int
	var ggpPrice *big.Int

	// Get data
	wg.Go(func() error {
		var err error
		remainingGgpBalance, err = auction.GetRemainingGGPBalance(rp, nil)
		return err
	})
	wg.Go(func() error {
		var err error
		lotMinimumEthValue, err = protocol.GetLotMinimumEthValue(rp, nil)
		return err
	})
	wg.Go(func() error {
		var err error
		ggpPrice, err = network.GetGGPPrice(rp, nil)
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return false, err
	}

	// Calculate lot minimum GGP amount
	var tmp big.Int
	var lotMinimumGgpAmount big.Int
	tmp.Mul(lotMinimumEthValue, eth.EthToWei(1))
	lotMinimumGgpAmount.Quo(&tmp, ggpPrice)

	// Return
	return (remainingGgpBalance.Cmp(&lotMinimumGgpAmount) >= 0), nil

}

// Get all lot count details
func getAllLotCountDetails(rp *rocketpool.RocketPool, bidderAddress common.Address) ([]lotCountDetails, error) {

	// Get lot count
	lotCount, err := auction.GetLotCount(rp, nil)
	if err != nil {
		return []lotCountDetails{}, err
	}

	// Load details in batches
	details := make([]lotCountDetails, lotCount)
	for bsi := uint64(0); bsi < lotCount; bsi += LotCountDetailsBatchSize {

		// Get batch start & end index
		lsi := bsi
		lei := bsi + LotCountDetailsBatchSize
		if lei > lotCount {
			lei = lotCount
		}

		// Load details
		var wg errgroup.Group
		for li := lsi; li < lei; li++ {
			li := li
			wg.Go(func() error {
				lotDetails, err := getLotCountDetails(rp, bidderAddress, li)
				if err == nil {
					details[li] = lotDetails
				}
				return err
			})
		}
		if err := wg.Wait(); err != nil {
			return []lotCountDetails{}, err
		}

	}

	// Return
	return details, nil

}

// Get a lot's count details
func getLotCountDetails(rp *rocketpool.RocketPool, bidderAddress common.Address, lotIndex uint64) (lotCountDetails, error) {

	// Data
	var wg errgroup.Group
	var addressBidAmount *big.Int
	var cleared bool
	var remainingGgp *big.Int
	var ggpRecovered bool

	// Get address bid amount
	wg.Go(func() error {
		var err error
		addressBidAmount, err = auction.GetLotAddressBidAmount(rp, lotIndex, bidderAddress, nil)
		return err
	})

	// Get lot cleared status
	wg.Go(func() error {
		var err error
		cleared, err = auction.GetLotIsCleared(rp, lotIndex, nil)
		return err
	})

	// Get lot remaining GGP amount
	wg.Go(func() error {
		var err error
		remainingGgp, err = auction.GetLotRemainingGGPAmount(rp, lotIndex, nil)
		return err
	})

	// Get lot GGP recovered status
	wg.Go(func() error {
		var err error
		ggpRecovered, err = auction.GetLotGGPRecovered(rp, lotIndex, nil)
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return lotCountDetails{}, err
	}

	// Return
	return lotCountDetails{
		AddressHasBid:   (addressBidAmount.Cmp(big.NewInt(0)) > 0),
		Cleared:         cleared,
		HasRemainingGgp: (remainingGgp.Cmp(big.NewInt(0)) > 0),
		GgpRecovered:    ggpRecovered,
	}, nil

}

// Get all lot details
func getAllLotDetails(rp *rocketpool.RocketPool, bidderAddress common.Address) ([]api.LotDetails, error) {

	// Get lot count
	lotCount, err := auction.GetLotCount(rp, nil)
	if err != nil {
		return []api.LotDetails{}, err
	}

	// Load details in batches
	details := make([]api.LotDetails, lotCount)
	for bsi := uint64(0); bsi < lotCount; bsi += LotDetailsBatchSize {

		// Get batch start & end index
		lsi := bsi
		lei := bsi + LotDetailsBatchSize
		if lei > lotCount {
			lei = lotCount
		}

		// Load details
		var wg errgroup.Group
		for li := lsi; li < lei; li++ {
			li := li
			wg.Go(func() error {
				lotDetails, err := getLotDetails(rp, bidderAddress, li)
				if err == nil {
					details[li] = lotDetails
				}
				return err
			})
		}
		if err := wg.Wait(); err != nil {
			return []api.LotDetails{}, err
		}

	}

	// Return
	return details, nil

}

// Get a lot's details
func getLotDetails(rp *rocketpool.RocketPool, bidderAddress common.Address, lotIndex uint64) (api.LotDetails, error) {

	// Get lot details
	details, err := auction.GetLotDetailsWithBids(rp, lotIndex, bidderAddress, nil)
	if err != nil {
		return api.LotDetails{}, err
	}

	// Check lot conditions
	addressHasBid := (details.AddressBidAmount.Cmp(big.NewInt(0)) > 0)
	hasRemainingGgp := (details.RemainingGGPAmount.Cmp(big.NewInt(0)) > 0)

	// Return
	return api.LotDetails{
		Details:              details,
		ClaimAvailable:       (addressHasBid && details.Cleared),
		BiddingAvailable:     (!details.Cleared && hasRemainingGgp),
		GGPRecoveryAvailable: (details.Cleared && hasRemainingGgp && !details.GGPRecovered),
	}, nil

}
