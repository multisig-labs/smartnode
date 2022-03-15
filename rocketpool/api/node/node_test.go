package node

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/go-homedir"
	"github.com/prysmaticlabs/prysm/v2/testing/require"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/smartnode/shared"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/rocket-pool/smartnode/shared/types/api"
	apitypes "github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func waitForTransaction(c *cli.Context, hash common.Hash) (*apitypes.APIResponse, error) {

	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := apitypes.APIResponse{}
	_, err = utils.WaitForTransaction(rp.Client, hash)
	if err != nil {
		return nil, err
	}

	// Return response
	return &response, nil

}

func initApp() (*cli.App, string, string) {
	app := cli.NewApp()
	app.Name = "gogopool"
	app.Usage = "GoGoPool CLI"
	app.Version = shared.RocketPoolVersion
	app.Authors = []cli.Author{
		{
			Name:  "Steven Gates",
			Email: "steven@multisiglabs.org",
		},
	}
	app.Copyright = "(c) 2022 Multisig Labs, Inc."

	// Get the config path from the arguments (or use the default)
	configPath := "~/.gogopool"
	for index, arg := range os.Args {
		if arg == "-c" || arg == "--config-path" {
			if len(os.Args)-1 == index {
				fmt.Fprintf(os.Stderr, "Expected config path after %s but none was given.\n", arg)
				os.Exit(1)
			}
			configPath = os.Args[index+1]
		}
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-path, c",
			Usage: "GoGoPool config asset `path`",
			Value: "~/.gogopool",
		},
	}

	// Get and parse the config file
	configFile := fmt.Sprintf("%s/%s", configPath, rocketpool.GlobalConfigFile)

	expandedConfigPath, err := homedir.Expand(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the global config file path: %s\n", err.Error())
		os.Exit(1)
	}
	// Stop if the config file doesn't exist yet
	_, err = os.Stat(expandedConfigPath)
	if !os.IsNotExist(err) {
		configBytes, err := ioutil.ReadFile(expandedConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load the global config file: %s\n", err.Error())
			os.Exit(1)
		}
		_, err = config.Parse(configBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse the global config file: %s\n", err.Error())
			os.Exit(1)
		}

	}

	settingsFile := fmt.Sprintf("%s/%s", "~/.gogopool", rocketpool.UserConfigFile)
	expandedSettingsPath, err := homedir.Expand(settingsFile)

	return app, expandedConfigPath, expandedSettingsPath
}

func TestCanNodeSend(t *testing.T) {
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)
	a := big.NewInt(0)
	a.SetString("998999343750000000001", 10)

	resp, err := canNodeSend(c, a, "avax")
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err, "node canNodeSend should not return error")
	fmt.Println(resp)

}

func TestCanNodeRegister(t *testing.T) {
	timezone := "Etc/UTC"
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)

	nodeResponse, err := canRegisterNode(c, timezone)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err, "node register should not return error")

	fmt.Println(nodeResponse)

}

func TestNodeRegister(t *testing.T) {
	timezone := "Etc/UTC"
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)

	nodeResponse, err := registerNode(c, timezone)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err, "node register should not return error")

	fmt.Println(nodeResponse)

}

func TestCanNodeStake(t *testing.T) {
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)
	stakeAmount := new(big.Int)
	stakeAmount.SetString("100", 10)
	nodeResponse, err := canNodeStakeRpl(c, stakeAmount)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err, "node register should not return error")

	fmt.Println(nodeResponse)

}

func prettyPrintResponse(response interface{}) {
	formattedResponse, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(formattedResponse))

}

func TestNodeStatus(t *testing.T) {
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)
	stakeAmount := new(big.Int)
	stakeAmount.SetString("100", 10)
	nodeResponse, err := getStatus(c)
	if err != nil {
		fmt.Println(err)
	}
	prettyPrintResponse(nodeResponse)
	assert.Nil(t, err, "node register should not return error")

}
func TestNodeStakeRPL(t *testing.T) {
	app, configPath, settingsPath := initApp()

	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)
	_, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}

	stakeAmount := new(big.Int)
	stakeAmount.SetString("100000000", 10)

	// Calculate max uint256 value
	maxApproval := big.NewInt(2)
	maxApproval = maxApproval.Exp(maxApproval, big.NewInt(256), nil)
	maxApproval = maxApproval.Sub(maxApproval, big.NewInt(1))

	// Approve RPL for staking
	response, err := approveRpl(c, maxApproval)
	if err != nil {
		fmt.Println(err)
	}
	hash := response.ApproveTxHash

	fmt.Printf("Approving RPL for staking...\n")

	if _, err = waitForTransaction(c, hash); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully approved staking access to RPL.")

	nodeResponse, err := stakeRpl(c, stakeAmount)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err, "node register should not return error")
	require.NoError(t, err)
	fmt.Println(nodeResponse)

}

const DefaultMaxNodeFeeSlippage = 0.01 // 1% below current network fee

func TestNodeDepositAVAX(t *testing.T) {
	app, configPath, settingsPath := initApp()

	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)
	_, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}

	stakeAmount := new(big.Int)
	stakeAmount.SetString("100000000", 10)

	// Calculate max uint256 value
	maxApproval := big.NewInt(2)
	maxApproval = maxApproval.Exp(maxApproval, big.NewInt(256), nil)
	maxApproval = maxApproval.Sub(maxApproval, big.NewInt(1))

	buffer := make([]byte, 32)
	_, err = rand.Read(buffer)
	if err != nil {
		fmt.Println(fmt.Errorf("Error generating random salt: %w", err))
	}
	var salt = big.NewInt(0).SetBytes(buffer)
	var minNodeFee float64

	nodeFees := api.NodeFeeResponse{}
	nodeFees.MaxNodeFee = 100
	nodeFees.MaxNodeFee = 20

	// Use default max slippage

	minNodeFee = nodeFees.NodeFee - DefaultMaxNodeFeeSlippage

	depositResponse, err := nodeDeposit(c, stakeAmount, minNodeFee, salt)
	if err != nil {
		return
	}

	prettyPrintResponse(depositResponse)

}

func TestContractCall(t *testing.T) {
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)

	// Get services
	if err := services.RequireNodeWallet(c); err != nil {
		fmt.Println(err)
	}
	if err := services.RequireRocketStorage(c); err != nil {
		fmt.Println(err)
	}
	_, err := services.GetWallet(c)
	if err != nil {
		fmt.Println(err)
	}

	//address := new(common.Address)
	//add := new(common.Address)
	//add.SetBytes([]byte("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"))
	//if err := ggp.GoGoStorageContract.Call(nil, address, "setGuardian", add); err != nil {
	//	fmt.Errorf("Could not get guardian address: %w", err)
	//}
	//fmt.Println(address)

	//contract, err := ggp.GetContract("rocketDAONodeTrustedSettingsMembers")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//getAndPrintBond(c)
	//setNewBondPrice(c)
	//getAndPrintBond(c)

	//getIsBootstrappedMode(c)
	//getContractAddress(c)
	genericContractCall(c)
}

func genericContractCall(c *cli.Context) {
	ggp, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}

	contract, err := ggp.GetContract("rocketDAOProtocolSettingsNode")
	if err != nil {
		fmt.Println(err)
	}

	//settingBool := new(bool)
	////address.SetBytes([]byte("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"))
	////add := new(common.Address)
	////add.SetBytes([]byte("0x3b7e31510e84988222f4a63db260d36c503d57d2"))
	//if err := contract.Call(nil, settingBool, "getSettingBool", "node.registration.enabled"); err != nil {
	//	fmt.Errorf("Could not get address: %w", err)
	//}
	//
	//fmt.Println(*settingBool)

	//settingBool := new(bool)
	//address.SetBytes([]byte("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"))
	add := new(common.Address)
	//add.SetBytes([]byte("0x3b7e31510e84988222f4a63db260d36c503d57d2"))
	//uintt := new(*big.Int)
	if err := contract.Call(nil, add, "getContractAddress", "rocketDAOProtocolSettingsNode"); err != nil {
		fmt.Errorf("Could not get address: %w", err)
	}

	fmt.Println(*add)
}

func getContractAddress(c *cli.Context) {
	ggp, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}

	contract, err := ggp.GetContract("rocketStorage")
	if err != nil {
		fmt.Println(err)
	}

	address := new(common.Address)
	//address.SetBytes([]byte("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"))
	//add := new(common.Address)
	//add.SetBytes([]byte("0x3b7e31510e84988222f4a63db260d36c503d57d2"))
	if err := contract.Call(nil, address, "getGuardian"); err != nil {
		fmt.Errorf("Could not get address: %w", err)
	}

	fmt.Println(address)
}

func getIsBootstrappedMode(c *cli.Context) {
	ggp, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}

	contract, err := ggp.GetContract("rocketDAONodeTrusted")
	if err != nil {
		fmt.Println(err)
	}

	bond := new(bool)
	if err := contract.Call(nil, bond, "getBootstrapModeDisabled"); err != nil {
		fmt.Errorf("Could not get bond amount: %w", err)
	}

	fmt.Println(*bond)
}

func setNewBondPrice(c *cli.Context) {
	ggp, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}
	contract, err := ggp.GetContract("rocketDAONodeTrustedSettingsMembers")

	newBond := big.NewInt(0)
	newBond.SetString("1000000000000000000000", 10)
	if err := contract.Call(nil, newBond, "setSettingUint", "members.rplbond", newBond); err != nil {
		fmt.Errorf("Could not get bond amount: %w", err)
	}

}

func getAndPrintBond(c *cli.Context) {
	ggp, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}

	contract, err := ggp.GetContract("rocketDAONodeTrustedSettingsMembers")
	if err != nil {
		fmt.Println(err)
	}

	bond := new(*big.Int)
	if err := contract.Call(nil, bond, "getRPLBond"); err != nil {
		fmt.Errorf("Could not get bond amount: %w", err)
	}

	fmt.Println(*bond)
}

//
//func TestCanNodeStakeGgp(t *testing.T) {
//	app, configPath, settingsPath := initApp()
//	set := flag.NewFlagSet("config-path", 0)
//	set.String("config", configPath, "doc")
//	set.String("settings", settingsPath, "doc")
//	c := cli.NewContext(app, set, nil)
//
//	canNodeStakeGgp(c, 100000)
//}
