package config

// Constants
const defaultPocketGatewayMainnet string = "lb/613bb4ae8c124d00353c40a1"
const defaultPocketGatewayPrater string = "lb/6126b4a783e49000343a3a47"

// Configuration for Pocket
type PocketConfig struct {
	// Common parameters that Pocket doesn't support and should be hidden
	UnsupportedCommonParams []string

	// Compatible consensus clients
	CompatibleConsensusClients []ConsensusClient

	// The Pocket gateway ID
	GatewayID Parameter
}

// Generates a new Pocket configuration
func NewPocketConfig(config *MasterConfig) *PocketConfig {
	return &PocketConfig{
		UnsupportedCommonParams: []string{ecWsPortID},

		CompatibleConsensusClients: []ConsensusClient{
			ConsensusClient_Lighthouse,
			ConsensusClient_Prysm,
			ConsensusClient_Teku,
		},

		GatewayID: Parameter{
			ID:          "gatewayID",
			Name:        "Gateway ID",
			Description: "If you would like to use a custom gateway for Pocket instead of the default Rocket Pool gateway, enter it here.",
			Type:        ParameterType_String,
			Default: map[Network]interface{}{
				Network_Mainnet: defaultPocketGatewayMainnet,
				Network_Prater:  defaultPocketGatewayPrater,
			},
			AffectsContainers:    []ContainerID{ContainerID_Eth1},
			EnvironmentVariables: []string{"POCKET_GATEWAY_ID"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Handle a network change on all of the parameters
func (config *PocketConfig) changeNetwork(oldNetwork Network, newNetwork Network) {
	changeNetworkForParameter(&config.GatewayID, oldNetwork, newNetwork)
}