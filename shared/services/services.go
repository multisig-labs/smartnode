package services

import (
	"fmt"
	"github.com/rocket-pool/smartnode/shared/services/beacon/avalanchego"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	uc "github.com/rocket-pool/rocketpool-go/utils/client"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services/beacon"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/contracts"
	"github.com/rocket-pool/smartnode/shared/services/passwords"
	"github.com/rocket-pool/smartnode/shared/services/wallet"
	lhkeystore "github.com/rocket-pool/smartnode/shared/services/wallet/keystore/lighthouse"
	nmkeystore "github.com/rocket-pool/smartnode/shared/services/wallet/keystore/nimbus"
	prkeystore "github.com/rocket-pool/smartnode/shared/services/wallet/keystore/prysm"
	tkkeystore "github.com/rocket-pool/smartnode/shared/services/wallet/keystore/teku"
)

// Config
const DockerAPIVersion = "1.40"

// Service instances & initializers
var (
	cfg             config.RocketPoolConfig
	passwordManager *passwords.PasswordManager
	nodeWallet      *wallet.Wallet
	ethClientProxy  *uc.EthClientProxy
	rocketPool      *rocketpool.RocketPool
	oneInchOracle   *contracts.OneInchOracle
	ggpFaucet       *contracts.GGPFaucet
	beaconClient    beacon.Client
	docker          *client.Client

	initCfg             sync.Once
	initPasswordManager sync.Once
	initNodeWallet      sync.Once
	initEthClientProxy  sync.Once
	initRocketPool      sync.Once
	initOneInchOracle   sync.Once
	initGgpFaucet       sync.Once
	initBeaconClient    sync.Once
	initDocker          sync.Once
)

//
// Service providers
//

func GetConfig(c *cli.Context) (config.RocketPoolConfig, error) {
	return getConfig(c)
}

func GetPasswordManager(c *cli.Context) (*passwords.PasswordManager, error) {
	cfg, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	return getPasswordManager(cfg), nil
}

func GetWallet(c *cli.Context) (*wallet.Wallet, error) {
	cfg, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	pm := getPasswordManager(cfg)
	return getWallet(cfg, pm)
}

func GetEthClientProxy(c *cli.Context) (*uc.EthClientProxy, error) {
	cfg, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	ec, err := getEthClientProxy(cfg)
	if err != nil {
		return nil, err
	}
	return ec, nil
}

func GetRocketPool(c *cli.Context) (*rocketpool.RocketPool, error) {
	cfg, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	ec, err := getEthClientProxy(cfg)
	if err != nil {
		return nil, err
	}

	return getRocketPool(cfg, ec)
}

func GetOneInchOracle(c *cli.Context) (*contracts.OneInchOracle, error) {
	cfg, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	ec, err := getEthClientProxy(cfg)
	if err != nil {
		return nil, err
	}
	return getOneInchOracle(cfg, ec)
}

func GetGgpFaucet(c *cli.Context) (*contracts.GGPFaucet, error) {
	cfg, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	ec, err := getEthClientProxy(cfg)
	if err != nil {
		return nil, err
	}
	return getGgpFaucet(cfg, ec)
}

func GetBeaconClient(c *cli.Context) (beacon.Client, error) {
	cfg, err := getConfig(c)
	if err != nil {
		return nil, err
	}
	return getBeaconClient(cfg)
}

func GetDocker(c *cli.Context) (*client.Client, error) {
	return getDocker()
}

//
// Service instance getters
//

func getConfig(c *cli.Context) (config.RocketPoolConfig, error) {
	var err error
	initCfg.Do(func() {
		cfg, err = config.Load(c)
	})
	return cfg, err
}

func getPasswordManager(cfg config.RocketPoolConfig) *passwords.PasswordManager {
	initPasswordManager.Do(func() {
		passwordManager = passwords.NewPasswordManager(os.ExpandEnv(cfg.Smartnode.PasswordPath))
	})
	return passwordManager
}

func getWallet(cfg config.RocketPoolConfig, pm *passwords.PasswordManager) (*wallet.Wallet, error) {
	var err error
	initNodeWallet.Do(func() {
		var maxFee *big.Int
		var maxPriorityFee *big.Int
		var gasLimit uint64
		maxFee, err = cfg.GetMaxFee()
		if err != nil {
			return
		}
		maxPriorityFee, err = cfg.GetMaxPriorityFee()
		if err != nil {
			return
		}
		gasLimit, err = cfg.GetGasLimit()
		if err != nil {
			return
		}
		nodeWallet, err = wallet.NewWallet(os.ExpandEnv(cfg.Smartnode.WalletPath), cfg.Chains.Platform.ChainID, maxFee, maxPriorityFee, gasLimit, pm)
		if err != nil {
			return
		}
		lighthouseKeystore := lhkeystore.NewKeystore(os.ExpandEnv(cfg.Smartnode.ValidatorKeychainPath), pm)
		nimbusKeystore := nmkeystore.NewKeystore(os.ExpandEnv(cfg.Smartnode.ValidatorKeychainPath), pm)
		prysmKeystore := prkeystore.NewKeystore(os.ExpandEnv(cfg.Smartnode.ValidatorKeychainPath), pm)
		tekuKeystore := tkkeystore.NewKeystore(os.ExpandEnv(cfg.Smartnode.ValidatorKeychainPath), pm)
		nodeWallet.AddKeystore("lighthouse", lighthouseKeystore)
		nodeWallet.AddKeystore("nimbus", nimbusKeystore)
		nodeWallet.AddKeystore("prysm", prysmKeystore)
		nodeWallet.AddKeystore("teku", tekuKeystore)
	})
	return nodeWallet, err
}

func getEthClientProxy(cfg config.RocketPoolConfig) (*uc.EthClientProxy, error) {
	var err error
	initEthClientProxy.Do(func() {
		reconnectDelay, err := time.ParseDuration(cfg.Chains.Platform.ReconnectDelay)
		if err != nil {
			return
		}
		if cfg.Chains.Platform.Client.Selected == "" {
			ethClientProxy = uc.NewEth1ClientProxy(reconnectDelay, cfg.Chains.Platform.Provider)
		} else {
			ethClientProxy = uc.NewEth1ClientProxy(reconnectDelay, cfg.Chains.Platform.Provider, cfg.Chains.Platform.FallbackProvider)
		}
	})
	return ethClientProxy, err
}

func getRocketPool(cfg config.RocketPoolConfig, client *uc.EthClientProxy) (*rocketpool.RocketPool, error) {
	var err error
	initRocketPool.Do(func() {
		rocketPool, err = rocketpool.NewRocketPool(client, common.HexToAddress(cfg.Rocketpool.StorageAddress))
	})
	return rocketPool, err
}

func getOneInchOracle(cfg config.RocketPoolConfig, client *uc.EthClientProxy) (*contracts.OneInchOracle, error) {
	var err error
	initOneInchOracle.Do(func() {
		oneInchOracle, err = contracts.NewOneInchOracle(common.HexToAddress(cfg.Rocketpool.OneInchOracleAddress), client)
	})
	return oneInchOracle, err
}

func getGgpFaucet(cfg config.RocketPoolConfig, client *uc.EthClientProxy) (*contracts.GGPFaucet, error) {
	var err error
	initGgpFaucet.Do(func() {
		ggpFaucet, err = contracts.NewGGPFaucet(common.HexToAddress(cfg.Rocketpool.GGPFaucetAddress), client)
	})
	return ggpFaucet, err
}

func getBeaconClient(cfg config.RocketPoolConfig) (beacon.Client, error) {
	var err error
	initBeaconClient.Do(func() {
		switch cfg.Chains.Platform.Client.Selected {
		//case "lighthouse":
		//	beaconClient = lighthouse.NewClient(cfg.Chains.Platform.Provider)
		//case "nimbus":
		//	beaconClient = nimbus.NewClient(cfg.Chains.Platform.Provider)
		//case "prysm":
		//	beaconClient = prysm.NewClient(cfg.Chains.Platform.Provider)
		//case "teku":
		//	beaconClient = teku.NewClient(cfg.Chains.Platform.Provider)
		case "avalanchego":
			beaconClient = avalanchego.NewClient(cfg.Chains.Platform.Provider)
		default:
			err = fmt.Errorf("Unknown Eth 2.0 client '%s' selected", cfg.Chains.Platform.Client.Selected)
		}
	})
	return beaconClient, err
}

func getDocker() (*client.Client, error) {
	var err error
	initDocker.Do(func() {
		docker, err = client.NewClientWithOpts(client.WithVersion(DockerAPIVersion))
	})
	return docker, err
}
