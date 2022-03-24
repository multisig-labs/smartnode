package watchtower

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/rocket-pool/smartnode/shared"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	"github.com/urfave/cli"
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

func prettyPrintResponse(response interface{}) {
	formattedResponse, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(formattedResponse))

}

func TestWithdrawAndStake(t *testing.T) {
	app, configPath, settingsPath := initApp()

	set := flag.NewFlagSet("config-path", 0)
	set.String("config", configPath, "doc")
	set.String("settings", settingsPath, "doc")
	c := cli.NewContext(app, set, nil)
	_, err := services.GetRocketPool(c)
	if err != nil {
		fmt.Println(err)
	}

	err = WithdrawAndStake(c)
	if err != nil {
		fmt.Println(err)
	}
	//prettyPrintResponse()

}
