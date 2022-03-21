package watchtower

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/urfave/cli"
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

const exportCChainParams = "{ \"to\":\"P-local192yta3e8v9j3em8lxa28w5pnj6m3ga9hqdygtj\", \"amount\": 2000, \"assetID\": \"AVAX\", \"username\":\"admin\", \"password\":\"adminadmin\" }"
const importPChainParams = "{ \"to\": \"P-local192yta3e8v9j3em8lxa28w5pnj6m3ga9hqdygtj\", \"from\": [\"P-local192yta3e8v9j3em8lxa28w5pnj6m3ga9hqdygtj\"], \"changeAddr\": \"P-local192yta3e8v9j3em8lxa28w5pnj6m3ga9hqdygtj\", \"username\": \"admin\", \"password\": \"adminadmin\" }"

const importKey = "{ \"username\" :\"admin\", \"password\":\"adminadmin\", \"privateKey\":\"PrivateKey-vTsBQDa1rD4XnS8akLz1Ba4uDDT3awfB8j9X15gQMg8yBRPFf\"}"

// Get the correct withdrawal credentials and pubkeys for each minipool
func stakeMinipools(c *cli.Context, rp *rocketpool.RocketPool, minipoolAddresses []common.Address) error {

	//ac, err := services.GetBeaconClient(c)
	//if err != nil {
	//	return err
	//}

	// ONE TIME import keys
	//cKeyImport, err := ac.MakeRPCCall("avax.importKey", "/ext/bc/C/avax", importKey)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(cKeyImport)
	//
	//pKeyImport, err := ac.MakeRPCCall("platform.importKey", "/ext/C", importKey)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(pKeyImport)
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

		// Export funds to P Chain address
		//exportResp, err := ac.MakeRPCCall("avax.export", "/ext/bc/C/avax", exportCChainParams)
		//if err != nil {
		//	return err
		//}
		//
		//fmt.Println(exportResp)
		//
		//wait for transaction?
		//importResp, err := ac.MakeRPCCall("platform.importAVAX", "/ext/P", importPChainParams)
		//if err != nil {
		//	return err
		//}
		//fmt.Println(importResp)

		// Move balance to P Chain address

		// Construct staking request object

		// Make RPC call to start staking
	}

	return nil

}
