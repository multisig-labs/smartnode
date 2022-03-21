package api

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/auction"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

type AuctionStatusResponse struct {
	Status              string   `json:"status"`
	Error               string   `json:"error"`
	TotalGGPBalance     *big.Int `json:"totalGGPBalance"`
	AllottedGGPBalance  *big.Int `json:"allottedGGPBalance"`
	RemainingGGPBalance *big.Int `json:"remainingGGPBalance"`
	CanCreateLot        bool     `json:"canCreateLot"`
	LotCounts           struct {
		ClaimAvailable       int `json:"claimAvailable"`
		BiddingAvailable     int `json:"biddingAvailable"`
		GGPRecoveryAvailable int `json:"ggpRecoveryAvailable"`
	} `json:"lotCounts"`
}

type AuctionLotsResponse struct {
	Status string       `json:"status"`
	Error  string       `json:"error"`
	Lots   []LotDetails `json:"lots"`
}
type LotDetails struct {
	Details              auction.LotDetails `json:"details"`
	ClaimAvailable       bool               `json:"claimAvailable"`
	BiddingAvailable     bool               `json:"biddingAvailable"`
	GGPRecoveryAvailable bool               `json:"ggpRecoveryAvailable"`
}

type CanCreateLotResponse struct {
	Status              string             `json:"status"`
	Error               string             `json:"error"`
	CanCreate           bool               `json:"canCreate"`
	InsufficientBalance bool               `json:"insufficientBalance"`
	CreateLotDisabled   bool               `json:"createLotDisabled"`
	GasInfo             rocketpool.GasInfo `json:"gasInfo"`
}
type CreateLotResponse struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	LotId  uint64      `json:"lotId"`
	TxHash common.Hash `json:"txHash"`
}

type CanBidOnLotResponse struct {
	Status           string             `json:"status"`
	Error            string             `json:"error"`
	CanBid           bool               `json:"canBid"`
	DoesNotExist     bool               `json:"doesNotExist"`
	BiddingEnded     bool               `json:"biddingEnded"`
	GGPExhausted     bool               `json:"ggpExhausted"`
	BidOnLotDisabled bool               `json:"bidOnLotDisabled"`
	GasInfo          rocketpool.GasInfo `json:"gasInfo"`
}
type BidOnLotResponse struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	TxHash common.Hash `json:"txHash"`
}

type CanClaimFromLotResponse struct {
	Status           string             `json:"status"`
	Error            string             `json:"error"`
	CanClaim         bool               `json:"canClaim"`
	DoesNotExist     bool               `json:"doesNotExist"`
	NoBidFromAddress bool               `json:"noBidFromAddress"`
	NotCleared       bool               `json:"notCleared"`
	GasInfo          rocketpool.GasInfo `json:"gasInfo"`
}
type ClaimFromLotResponse struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	TxHash common.Hash `json:"txHash"`
}

type CanRecoverGGPFromLotResponse struct {
	Status              string             `json:"status"`
	Error               string             `json:"error"`
	CanRecover          bool               `json:"canRecover"`
	DoesNotExist        bool               `json:"doesNotExist"`
	BiddingNotEnded     bool               `json:"biddingNotEnded"`
	NoUnclaimedGGP      bool               `json:"noUnclaimedGgp"`
	GGPAlreadyRecovered bool               `json:"ggpAlreadyRecovered"`
	GasInfo             rocketpool.GasInfo `json:"gasInfo"`
}
type RecoverGGPFromLotResponse struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	TxHash common.Hash `json:"txHash"`
}
