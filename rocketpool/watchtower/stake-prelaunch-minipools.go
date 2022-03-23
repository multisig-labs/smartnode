package watchtower

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/urfave/cli"
	"time"
)

// Submit scrub minipools
func WithdrawAndStake(c *cli.Context) error {

	// Get services
	if err := services.RequireNodeRegistered(c); err != nil {
		return err
	}
	//w, err := services.GetWallet(c)
	//if err != nil {
	//	return err
	//}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return err
	}
	// Get node account
	//nodeAccount, err := w.GetNodeAccount()
	//if err != nil {
	//	return err
	//}

	// Get trusted node status
	//nodeTrusted, err := trustednode.GetMemberExists(rp, nodeAccount.Address, nil)
	//if err != nil {
	//	return err
	//}
	//if !(nodeTrusted) {
	//	return nil
	//}

	// Log
	fmt.Println("Checking for minipools with prelaunch status to stake...")

	// Get minipools in prelaunch status
	minipoolAddresses, err := minipool.GetPrelaunchMinipoolAddresses(rp, nil)
	if err != nil {
		return err
	}
	totalMinipools := len(minipoolAddresses)
	if totalMinipools == 0 {
		fmt.Println("No minipools in prelaunch.")
		return nil
	}

	// Get the correct withdrawal credentials and validator pubkeys for each minipool
	err = stakeMinipools(c, rp, minipoolAddresses)
	if err != nil {
		return err
	}
	return nil

}

const exportCChainParams = "{ \"to\":\"P-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\", \"assetID\": \"AVAX\", \"amount\": 7000000000, \"username\":\"admin\", \"password\":\"Lolsaldkfjxckmnvipop123!@#\" }"
const importPChainParams = "{ \"username\":\"admin\", \"password\":\"Lolsaldkfjxckmnvipop123!@#\", \"sourceChain\": \"C\", \"to\":\"P-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\" }"

var stakePChainParams = "{ \"nodeID\":\"NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg\", \"rewardAddress\":\"P-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\", \"from\": [\"P-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\"], \"changeAddr\": \"P-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\", \"startTime\":'" + fmt.Sprintf("%d", time.Now().Unix()+60) + "', \"endTime\":'" + fmt.Sprintf("%d", time.Now().AddDate(0, 0, 14).Unix()) + "', \"stakeAmount\":7000000000, \"delegationFeeRate\":10, \"username\":\"admin\", \"password\":\"Lolsaldkfjxckmnvipop123!@#\" }"

//const importXChainParams = "{ \"to\":\"X-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\", \"sourceChain\":\"C\", \"username\":\"admin\", \"password\":\"Lolsaldkfjxckmnvipop123!@#\" }"
//const exportXChainParams = "{ \"to\":\"P-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\", \"amount\": 1000000000000, \"assetID\": \"AVAX\", \"from\":[\"X-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\"], \"changeAddr\":\"X-local12hnt0379l7vpfxryyf8guwgh7afyqfm0kwhm7u\", \"username\":\"admin\", \"password\":\"Lolsaldkfjxckmnvipop123!@#\" }"

// Get the correct withdrawal credentials and pubkeys for each minipool
func stakeMinipools(c *cli.Context, rp *rocketpool.RocketPool, minipoolAddresses []common.Address) error {

	ac, err := services.GetBeaconClient(c)
	if err != nil {
		return err
	}

	for _, minipoolAddress := range minipoolAddresses {
		// Create a minipool contract wrapper for the given address
		mp, err := minipool.NewMinipool(rp, minipoolAddress)
		if err != nil {
			fmt.Printf("Error creating minipool wrapper for %s: %s", minipoolAddress.Hex(), err.Error())
			continue
		}
		fmt.Println(mp.GetNodeDepositBalance(nil))
		fmt.Println(mp.GetNodeId(nil))
		fmt.Println(mp.Address)
		fmt.Println(mp.GetBalance(nil))

		// Get addressess

		// TODO derive addresses from the users node wallet. this is currently hardcoded to the wallet address for now.
		//pChainAddress := common.HexToAddress("P-local192yta3e8v9j3em8lxa28w5pnj6m3ga9hqdygtj")
		//xChainAdress := common.HexToAddress("X-local192yta3e8v9j3em8lxa28w5pnj6m3ga9hqdygtj")
		//cChainAddress := common.HexToAddress("0x3b7e31510e84988222f4a63db260d36c503d57d2")

		// Withdraw balance to this C Chain address
		//err = mp.WithdrawBalanceForStaking(nil)
		//if err != nil {
		//	return err
		//}

		// Export funds to X Chain address
		exportResp, err := ac.MakeRPCCall("avax.exportAVAX", "/ext/bc/C/avax", exportCChainParams)
		if err != nil {
			return err
		}

		fmt.Println(exportResp)
		time.Sleep(5 * time.Second)

		// Import to P Chain address
		importPResp, err := ac.MakeRPCCall("platform.importAVAX", "/ext/P", importPChainParams)
		if err != nil {
			return err
		}
		fmt.Println(importPResp)
		time.Sleep(5 * time.Second)
		time.Sleep(5)
		stakePChainResp, err := ac.MakeRPCCall("platform.addValidator", "/ext/P", importPChainParams)
		if err != nil {
			return err
		}
		fmt.Println(stakePChainResp)
		// Make RPC call to start staking
	}

	return nil

}
