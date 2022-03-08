package service

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

// Configure the Rocket Pool service
func configureService(c *cli.Context) error {

	colorReset := "\033[0m"
	colorYellow := "\033[33m"

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Load configs
	globalConfig, err := rp.LoadGlobalConfig()
	if err != nil {
		return err
	}
	userConfig, err := rp.LoadUserConfig()
	if err != nil {
		return err
	}

	showAdvanced := c.Bool("advanced")

	// Ask about advanced mode
	if !showAdvanced {
		if cliutils.Confirm("Some settings (such as port selection) come with recommended defaults.\n" +
			"Would you like to use them automatically? You can review them at the end of this setup.") {
			showAdvanced = false
		} else {
			showAdvanced = true
		}
	}

	// Configure eth1
	if err := configureChain(&(globalconfig.Chains.Platform), &(userconfig.Chains.Platform), "Eth 1.0", false, false, []string{}, true, showAdvanced); err != nil {
		return err
	}

	// Get the list of fallback eth1 clients
	var fallbackEth1Clients []string
	for _, eth1Client := range globalconfig.Chains.Platform.Client.Options {
		if eth1Client.Fallback {
			fallbackEth1Clients = append(fallbackEth1Clients, eth1Client.ID)
		}
	}

	// Get the list of compatible eth2 clients for the primary
	var compatibleEth2Clients []string
	compatibleString := globalconfig.Chains.Platform.GetSelectedClient().CompatibleEth2Clients
	if compatibleString != "" {
		compatibleEth2Clients = strings.Split(globalconfig.Chains.Platform.GetSelectedClient().CompatibleEth2Clients, ";")
	}

	// Configure eth1 fallback
	if cliutils.Confirm("Would you like to configure a second Eth 1.0 client to act as a fallback in case your primary Eth 1.0 client is unavailable?") {
		if err := configureChain(&(globalconfig.Chains.Platform), &(userconfig.Chains.PlatformFallback), "Eth 1.0 Fallback", false, true, fallbackEth1Clients, true, showAdvanced); err != nil {
			return err
		}

		// Get the list of compatible eth2 clients for the fallback
		var fallbackCompatibleEth2Clients []string
		fallbackCompatibleString := globalconfig.Chains.Platform.GetClientById(userconfig.Chains.PlatformFallback.Client.Selected).CompatibleEth2Clients
		if fallbackCompatibleString != "" {
			fallbackCompatibleEth2Clients = strings.Split(globalconfig.Chains.Platform.GetSelectedClient().CompatibleEth2Clients, ";")
		}

		if len(fallbackCompatibleEth2Clients) == 0 {
			// If the fallback is compatible with everything, do nothing
		} else if len(compatibleEth2Clients) == 0 {
			// If the primary client is compatible with all clients, defer to the fallback settings
			compatibleEth2Clients = fallbackCompatibleEth2Clients
		} else {
			// Remove any clients from the primary that aren't also in the fallback
			for i, compatibleClient := range compatibleEth2Clients {
				clientExists := false
				// See if this is also in the compatible fallback list
				for _, fallbackCompatibleClient := range fallbackCompatibleEth2Clients {
					if fallbackCompatibleClient == compatibleClient {
						clientExists = true
						break
					}
				}
				// Remove the client from the main compatible list
				if !clientExists {
					compatibleEth2Clients = append(compatibleEth2Clients[:i], compatibleEth2Clients[i+1:]...)
				}
			}
		}
	} else {
		userconfig.Chains.PlatformFallback = config.Chain{}
	}

	// Configure eth2
	if err := configureChain(&(globalConfig.Chains.Platform), &(userConfig.Chains.Platform), "Eth 2.0", true, false, compatibleEth2Clients, true, showAdvanced); err != nil {
		return err
	}

	// Configure metrics
	if err := configureMetrics(&(globalConfig.Metrics), &(userConfig.Metrics), showAdvanced); err != nil {
		return err
	}

	// Save user config
	if err := rp.SaveUserConfig(userConfig); err != nil {
		return err
	}

	// Print settings
	fmt.Println("=== ETH1 Settings ===")
	eth1Client := globalconfig.Chains.Platform.GetClientById(userconfig.Chains.Platform.Client.Selected)
	fmt.Printf("Selected client: %s\n", eth1Client.Name)
	for _, param := range userconfig.Chains.Platform.Client.Params {
		globalParam := eth1Client.GetParamByEnvName(param.Env)
		if globalParam != nil {
			fmt.Printf("%s: %s\n", globalParam.Name, param.Value)
		}
	}
	fmt.Println()

	if userconfig.Chains.PlatformFallback.Client.Selected != "" {
		fmt.Println("=== ETH1 Fallback Settings ===")
		eth1FallbackClient := globalconfig.Chains.Platform.GetClientById(userconfig.Chains.PlatformFallback.Client.Selected)
		fmt.Printf("Selected client: %s\n", eth1FallbackClient.Name)
		for _, param := range userconfig.Chains.PlatformFallback.Client.Params {
			globalParam := eth1FallbackClient.GetParamByEnvName(param.Env)
			if globalParam != nil {
				fmt.Printf("%s: %s\n", globalParam.Name, param.Value)
			}
		}
		fmt.Println()
	}

	fmt.Println("=== ETH2 Settings ===")
	eth2Client := globalConfig.Chains.Platform.GetClientById(userConfig.Chains.Platform.Client.Selected)
	fmt.Printf("Selected client: %s\n", eth2Client.Name)
	for _, param := range userConfig.Chains.Platform.Client.Params {
		globalParam := eth2Client.GetParamByEnvName(param.Env)
		if globalParam != nil {
			fmt.Printf("%s: %s\n", globalParam.Name, param.Value)
		}
	}
	fmt.Println()

	if userConfig.Metrics.Enabled {
		fmt.Println("=== Metrics Settings ===")
		for _, param := range userConfig.Metrics.Settings {
			fmt.Printf("%s: %s\n", globalConfig.Metrics.GetParamByEnvName(param.Env).Name, param.Value)
		}
		fmt.Println()
	}

	// Log & return
	fmt.Println("Done!\n")
	fmt.Printf("%sNOTE:\n", colorYellow)
	fmt.Printf("Please run 'rocketpool service stop' and 'rocketpool service start' to apply any changes you made.%s\n", colorReset)
	return nil

}

func configureMetrics(globalMetrics, userMetrics *config.Metrics, showAdvanced bool) error {

	// Prompt for enabling status
	enabled := cliutils.Confirm("Would you like to enable Rocket Pool's metrics dashboard?")
	userMetrics.Enabled = enabled

	// Prompt for params
	params := []config.UserParam{}
	if enabled {
		for _, param := range globalMetrics.Params {

			// Skip advanced parameters if they're disabled
			if param.Advanced && !showAdvanced {
				params = append(params, config.UserParam{
					Env:   param.Env,
					Value: param.Default,
				})
				continue
			}

			// Get expected param format
			var expectedFormat string
			if param.Regex != "" {
				expectedFormat = param.Regex
			} else if param.Required {
				expectedFormat = "^.+$"
			} else {
				expectedFormat = "^.*$"
			}

			// Get param label
			paramText := param.Name
			if !param.Required {
				blankText := "none"
				if param.BlankText != "" {
					blankText = param.BlankText
				}
				paramText += fmt.Sprintf(" (leave blank for %s)", blankText)
			}
			if param.Desc != "" {
				paramText += fmt.Sprintf("\n(%s)", param.Desc)
			}

			// Prompt for value
			var value string
			for {
				var isValid bool
				if param.Type == "bool" {
					value = cliutils.Prompt(fmt.Sprintf("Would you like to enable %s? [y/n]\n(%s)", param.Name, param.Desc), expectedFormat, fmt.Sprintf("Invalid %s", param.Name))
					isValid = true
				} else {
					value = cliutils.Prompt(fmt.Sprintf("Please enter the %s", paramText), expectedFormat, fmt.Sprintf("Invalid %s", param.Name))
					isValid = true
				}

				// Allow blanks for optional params
				if !param.Required && value == "" {
					value = param.Default
					break
				}

				// Type checking
				switch param.Type {
				case "uint":
					if _, err := strconv.ParseUint(value, 0, 0); err != nil {
						fmt.Printf("'%s' is not a valid value for %s, try again.\n", value, param.Name)
						isValid = false
					}
				case "uint16":
					if _, err := strconv.ParseUint(value, 0, 16); err != nil {
						fmt.Printf("'%s' is not a valid value for %s, try again.\n", value, param.Name)
						isValid = false
					}
				}

				// Continue if input is valid
				if isValid {
					break
				}

			}

			// Add param
			params = append(params, config.UserParam{
				Env:   param.Env,
				Value: value,
			})

		}
	}

	// Set unselected client params to blank strings to prevent docker-compose warnings
	for _, param := range globalMetrics.Params {

		// Cancel if param already set in selected client
		paramSet := false
		for _, userParam := range params {
			if param.Env == userParam.Env {
				paramSet = true
				break
			}
		}
		if paramSet {
			continue
		}

		// Add param
		params = append(params, config.UserParam{
			Env:   param.Env,
			Value: "",
		})

	}

	// Set config params
	userMetrics.Settings = params

	// Return
	return nil
}

// Configure a chain
func configureChain(globalChain, userChain *config.Chain, chainName string, defaultRandomClient bool, fallbackOnly bool, compatibleClients []string, includeSupermajority bool, showAdvanced bool) error {

	// Check client options
	if len(globalChain.Client.Options) == 0 {
		return fmt.Errorf("There are no available %s client options", chainName)
	}

	// Check if there's an existing client choice
	reuseClient := false
	if userChain.Client.Selected != "" {
		client := globalChain.GetClientById(userChain.Client.Selected)
		if client != nil {
			if (client.Supermajority && !includeSupermajority) || (!client.Fallback && fallbackOnly) {
				// If supermajorities are excluded and they're already using one, ignore it.
				// If in fallback clients only mode and this isn't one, ignore it.
				reuseClient = false
			} else {
				// Check the compatible clients and make sure the selected one is compatible
				isCompatible := false
				if len(compatibleClients) == 0 {
					isCompatible = true
				} else {
					for _, compatibleClient := range compatibleClients {
						if compatibleClient == client.ID {
							isCompatible = true
							break
						}
					}
				}
				if isCompatible {
					reuseClient = cliutils.Confirm(fmt.Sprintf(
						"Detected an existing %s client choice of %s.\nWould you like to continue using it?", chainName, client.Name))
				} else {
					fmt.Printf("%sWarning: your existing eth2 client [%s] is incompatible with your eth1 and fallback choices.%s\n\n", colorYellow, client.Name, colorReset)
				}
			}
		}
	}

	// If the user wants to pick a different client
	if !reuseClient {

		// Prompt for random client selection
		var randomClient bool
		if defaultRandomClient {
			randomClient = cliutils.Confirm(fmt.Sprintf("Would you like to run a random %s client (recommended)?", chainName))
		}

		// Create compatible client list
		var compatibleClientIndices []int
		if len(compatibleClients) > 0 {
			// Go through each client
			for clientIndex, clientId := range globalChain.Client.Options {
				// Go through the list of compatible clients
				for _, compatibleId := range compatibleClients {
					// If this client is compatible, add its index to the list of good ones
					if clientId.ID == compatibleId {
						compatibleClientIndices = append(compatibleClientIndices, clientIndex)
					}
				}
			}

			// Panic if the list is empty!
			if len(compatibleClientIndices) == 0 {
				return fmt.Errorf("There are no compatible %s clients available.", chainName)
			}
		} else {
			// Create an array with all of the client indices
			for i := range globalChain.Client.Options {
				compatibleClientIndices = append(compatibleClientIndices, i)
			}
		}

		// Select client
		var selected int
		if randomClient {
			// Make a new list indices in compatibleClientIndices which represent eligible clients
			var eligible []int
			for clientIndicesIndex, clientIndex := range compatibleClientIndices {
				if !includeSupermajority && globalChain.Client.Options[clientIndex].Supermajority {
					continue
				}

				if fallbackOnly && !globalChain.Client.Options[clientIndex].Fallback {
					continue
				}

				eligible = append(eligible, clientIndicesIndex)
			}

			if len(eligible) == 0 {
				// No clients could be selected, due to fallback/supermajority preferences
				return fmt.Errorf("There are no compatible %s clients available with Fallback %t and includeSupermajority %t", chainName, fallbackOnly, includeSupermajority)
			}

			rand.Seed(time.Now().UnixNano())
			selected = eligible[rand.Intn(len(eligible))]

		} else {
			clientOptions := make([]string, len(compatibleClientIndices))
			for oi, optionIndex := range compatibleClientIndices {
				option := globalChain.Client.Options[optionIndex]
				optionText := option.Name
				if option.Desc != "" {
					optionText += fmt.Sprintf(" %s\n\t\t%s\n", option.Desc, option.Link)
				}
				clientOptions[oi] = optionText
			}

			// Print incompatible clients
			var incompatibleClientNames []string
			if len(compatibleClients) > 0 {
				for _, clientId := range globalChain.Client.Options {
					incompatible := true
					for _, compatibleId := range compatibleClients {
						if clientId.ID == compatibleId {
							incompatible = false
							break
						}
					}
					if incompatible {
						incompatibleClientNames = append(incompatibleClientNames, clientId.Name)
					}
				}
				if len(incompatibleClientNames) > 0 {
					colorReset := "\033[0m"
					colorYellow := "\033[33m"
					fmt.Printf("%sIncompatible %s clients: %s\n\n%s", colorYellow, chainName, incompatibleClientNames, colorReset)
				}
			}

			selected, _ = cliutils.Select(fmt.Sprintf("Which %s client would you like to run?", chainName), clientOptions)
		}

		// Set selected client
		selectedIndex := compatibleClientIndices[selected]
		selectedId := globalChain.Client.Options[selectedIndex].ID
		globalChain.Client.Selected = selectedId
		userChain.Client.Selected = selectedId
	} else {
		globalChain.Client.Selected = userChain.Client.Selected
	}

	// Log
	fmt.Printf("%s %s client selected.\n", globalChain.GetSelectedClient().Name, chainName)
	fmt.Println("")

	// Warn about supermajority clients
	if globalChain.GetSelectedClient().Supermajority {
		fmt.Println("NOTE: Due to client diversity imbalance, we recommend you select a client with lower representation.")
		fmt.Println("Please visit https://docs.rocketpool.net/guides/node/eth-clients.html to learn more.")
		fmt.Println()
		useDifferent := cliutils.Confirm("Would you like to use a different client (recommended)?")
		if useDifferent {
			return configureChain(globalChain, userChain, chainName, defaultRandomClient, fallbackOnly, compatibleClients, false, showAdvanced)
		}
	}

	// Prompt for params
	params := []config.UserParam{}
	for _, param := range globalChain.GetSelectedClient().Params {

		// Skip advanced parameters if they're disabled
		if param.Advanced && !showAdvanced {
			params = append(params, config.UserParam{
				Env:   param.Env,
				Value: param.Default,
			})
			continue
		}

		// Get expected param format
		var expectedFormat string
		if param.Regex != "" {
			expectedFormat = param.Regex
		} else if param.Required {
			expectedFormat = "^.+$"
		} else {
			expectedFormat = "^.*$"
		}

		// Get param label
		paramText := param.Name
		if !param.Required {
			blankText := "none"
			if param.BlankText != "" {
				blankText = param.BlankText
			}
			paramText += fmt.Sprintf(" (leave blank for %s)", blankText)
		}
		if param.Desc != "" {
			paramText += fmt.Sprintf("\n(%s)", param.Desc)
		}

		// Prompt for value
		var value string
		for {
			var isValid bool
			if param.Type == "bool" {
				value = cliutils.Prompt(fmt.Sprintf("Would you like to enable %s? [y/n]\n(%s)", param.Name, param.Desc), expectedFormat, fmt.Sprintf("Invalid %s", param.Name))
				isValid = true
			} else {
				value = cliutils.Prompt(fmt.Sprintf("Please enter the %s", paramText), expectedFormat, fmt.Sprintf("Invalid %s", param.Name))
				isValid = true
			}

			// Allow blanks for optional params
			if !param.Required && value == "" {
				value = param.Default
				break
			}

			// Type checking
			switch param.Type {
			case "uint":
				if _, err := strconv.ParseUint(value, 0, 0); err != nil {
					fmt.Printf("'%s' is not a valid value for %s, try again.\n", value, param.Name)
					isValid = false
				}
			case "uint16":
				if _, err := strconv.ParseUint(value, 0, 16); err != nil {
					fmt.Printf("'%s' is not a valid value for %s, try again.\n", value, param.Name)
					isValid = false
				}
			}

			// Continue if input is valid
			if isValid {
				break
			}

		}

		// Add param
		params = append(params, config.UserParam{
			Env:   param.Env,
			Value: value,
		})

	}

	// Set unselected client params to blank strings to prevent docker-compose warnings
	for _, option := range globalChain.Client.Options {
		if option.ID == globalChain.Client.Selected {
			continue
		}
		for _, param := range option.Params {

			// Cancel if param already set in selected client
			paramSet := false
			for _, userParam := range params {
				if param.Env == userParam.Env {
					paramSet = true
					break
				}
			}
			if paramSet {
				continue
			}

			// Add param
			params = append(params, config.UserParam{
				Env:   param.Env,
				Value: "",
			})

		}
	}

	// Set config params
	userChain.Client.Params = params

	// Return
	return nil

}
