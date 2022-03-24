package node

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/tokens"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
)

func canNodeBurn(c *cli.Context, amountWei *big.Int, token string) (*api.CanNodeBurnResponse, error) {

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
	response := api.CanNodeBurnResponse{}

	// Get node account
	nodeAccount, err := w.GetNodeAccount()
	if err != nil {
		return nil, err
	}

	// Sync
	var wg errgroup.Group

	// Check node balance
	wg.Go(func() error {
		switch token {
		case "ggpavax":

			// Check node ggpAVAX balance
			ggpavaxBalanceWei, err := tokens.GetGGPAVAXBalance(rp, nodeAccount.Address, nil)
			if err != nil {
				return err
			}
			response.InsufficientBalance = (amountWei.Cmp(ggpavaxBalanceWei) > 0)

		}
		return nil
	})

	// Check token contract collateral
	wg.Go(func() error {
		switch token {
		case "ggpavax":

			// Check ggpAVAX collateral
			ggpavaxTotalCollateral, err := tokens.GetGGPAVAXTotalCollateral(rp, nil)
			if err != nil {
				return err
			}
			response.InsufficientCollateral = (amountWei.Cmp(ggpavaxTotalCollateral) > 0)

		}
		return nil
	})

	// Get gas estimate
	wg.Go(func() error {
		opts, err := w.GetNodeAccountTransactor()
		if err != nil {
			return err
		}
		switch token {
		case "ggpavax":
			gasInfo, err := tokens.EstimateBurnGGPAVAXGas(rp, amountWei, opts)
			if err == nil {
				response.GasInfo = gasInfo
			}
			return err
		}
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return nil, err
	}

	// Update & return response
	response.CanBurn = !(response.InsufficientBalance || response.InsufficientCollateral)
	return &response, nil

}

func nodeBurn(c *cli.Context, amountWei *big.Int, token string) (*api.NodeBurnResponse, error) {

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
	response := api.NodeBurnResponse{}

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

	// Handle token type
	switch token {
	case "ggpavax":

		// Burn ggpAVAX
		hash, err := tokens.BurnGGPAVAX(rp, amountWei, opts)
		if err != nil {
			return nil, err
		}
		response.TxHash = hash

	}

	// Return response
	return &response, nil

}
