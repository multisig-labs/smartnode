package network

import (
	"fmt"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)

func getGgpPrice(c *cli.Context) error {

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Get GGP price
	response, err := rp.GgpPrice()
	if err != nil {
		return err
	}

	// Print & return
	fmt.Printf("The current network GGP price is %.6f ETH.\n", math.RoundDown(eth.WeiToEth(response.GgpPrice), 6))
	fmt.Printf("Prices last updated at block: %d\n", response.GgpPriceBlock)
	return nil

}
