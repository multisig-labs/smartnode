package wallet

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/rocket-pool/smartnode/shared"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

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

func initConfig() {

}

func TestWalletInit(t *testing.T) {
	app, _, _ := initApp()
	err := app.Run([]string{"rocketpool", "wallet", "init"})
	assert.Nil(t, err, "wallet init should not return error")
}

func TestWalletStatusWhenNotInitialized(t *testing.T) {
	app, configPath, _ := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", "/Users/pisti/.gogopool/settings.yml", "doc")
	c := cli.NewContext(app, set, nil)

	resp, err := getStatus(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Nil(t, err, "wallet status should not return error")
	assert.Falsef(t, resp.WalletInitialized, "wallet should not be initialized")
	assert.Falsef(t, resp.PasswordSet, "wallet password should not be set")

}

func TestSetPassword(t *testing.T) {
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")

	c := cli.NewContext(app, set, nil)

	password := "123456789123"
	passwordResponse, err := setPassword(c, password)
	if err != nil {
		fmt.Println(err)
	}

	assert.Nil(t, err, "setpassword must not return error")
	assert.Truef(t, passwordResponse.Status == "success", "response was not success")

	// clean up
	pwPath, _ := homedir.Expand("~/.gogopool/data/password")
	err = os.Remove(pwPath)
	if err != nil {
		log.Fatal(err)
	}

}

func initializeWallet(c *cli.Context) (api.InitWalletResponse, error) {
	password := "123456789123"
	passwordResponse, err := setPassword(c, password)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(passwordResponse)

	// Initialize wallet
	walletResp, err := initWallet(c)
	fmt.Println(walletResp)
	if err != nil {
		fmt.Println(err)
		return api.InitWalletResponse{}, err
	}

	return api.InitWalletResponse{
		Status:         walletResp.Status,
		Error:          walletResp.Error,
		Mnemonic:       walletResp.Mnemonic,
		AccountAddress: walletResp.AccountAddress,
	}, nil
}

func removeWalletAndPassword() {
	pwPath, _ := homedir.Expand("~/.gogopool/data/password")
	err := os.Remove(pwPath)
	if err != nil {
		log.Fatal(err)
	}
	waPath, _ := homedir.Expand("~/.gogopool/data/wallet")
	err = os.Remove(waPath)
	if err != nil {
		log.Fatal(err)
	}
}
func TestWalletStatusWhenInitialized(t *testing.T) {
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)

	_, err := initializeWallet(c)
	assert.Nil(t, err, "wallet init should not return error")

	resp, err := getStatus(c)
	if err != nil {
		fmt.Println(err)
	}

	assert.Nil(t, err, "wallet status should not return error")
	assert.Truef(t, resp.WalletInitialized, "wallet should be initialized")
	assert.Truef(t, resp.PasswordSet, "wallet password should be set")

	// clean up
	//removeWalletAndPassword()
}

func TestWalletExport(t *testing.T) {
	removeWalletAndPassword()
	app, configPath, settingsPath := initApp()
	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)
	_, err := initializeWallet(c)
	//assert.Nil(t, err, "wallet init should not return error")

	wR, err := exportWallet(c)
	assert.Nil(t, err, "wallet export should not return error")
	fmt.Println(wR.AccountPrivateKey)
	fmt.Println(wR.Password)
	fmt.Println(wR.Wallet)

}
