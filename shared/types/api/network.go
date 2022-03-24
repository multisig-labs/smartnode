package api

import (
	"math/big"
)

type NodeFeeResponse struct {
	Status        string  `json:"status"`
	Error         string  `json:"error"`
	NodeFee       float64 `json:"nodeFee"`
	MinNodeFee    float64 `json:"minNodeFee"`
	TargetNodeFee float64 `json:"targetNodeFee"`
	MaxNodeFee    float64 `json:"maxNodeFee"`
}

type GgpPriceResponse struct {
	Status                 string   `json:"status"`
	Error                  string   `json:"error"`
	GgpPrice               *big.Int `json:"ggpPrice"`
	GgpPriceBlock          uint64   `json:"ggpPriceBlock"`
	MinPerMinipoolGgpStake *big.Int `json:"minPerMinipoolGgpStake"`
	MaxPerMinipoolGgpStake *big.Int `json:"maxPerMinipoolGgpStake"`
}

type NetworkStatsResponse struct {
	Status                    string  `json:"status"`
	Error                     string  `json:"error"`
	TotalValueLocked          float64 `json:"totalValueLocked"`
	DepositPoolBalance        float64 `json:"depositPoolBalance"`
	MinipoolCapacity          float64 `json:"minipoolCapacity"`
	StakerUtilization         float64 `json:"stakerUtilization"`
	NodeFee                   float64 `json:"nodeFee"`
	NodeCount                 uint64  `json:"nodeCount"`
	InitializedMinipoolCount  uint64  `json:"initializedMinipoolCount"`
	PrelaunchMinipoolCount    uint64  `json:"prelaunchMinipoolCount"`
	StakingMinipoolCount      uint64  `json:"stakingMinipoolCount"`
	WithdrawableMinipoolCount uint64  `json:"withdrawableMinipoolCount"`
	DissolvedMinipoolCount    uint64  `json:"dissolvedMinipoolCount"`
	FinalizedMinipoolCount    uint64  `json:"finalizedMinipoolCount"`
	GgpPrice                  float64 `json:"ggpPrice"`
	TotalGgpStaked            float64 `json:"totalGgpStaked"`
	EffectiveGgpStaked        float64 `json:"effectiveGgpStaked"`
	GgpavaxPrice              float64 `json:"ggpavaxPrice"`
}

type NetworkTimezonesResponse struct {
	Status         string            `json:"status"`
	Error          string            `json:"error"`
	TimezoneCounts map[string]uint64 `json:"timezoneCounts"`
	TimezoneTotal  uint64            `json:"timezoneTotal"`
	NodeTotal      uint64            `json:"nodeTotal"`
}
