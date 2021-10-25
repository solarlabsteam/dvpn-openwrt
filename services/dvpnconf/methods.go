package dvpnconf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/solarlabsteam/dvpn-openwrt/services/node"
	"github.com/solarlabsteam/dvpn-openwrt/utilities/appconf"

	"github.com/pelletier/go-toml"
)

const (
	BackEndTypeTest               = "test"
	DefaultRPCAddress             = "https://rpc.sentinel.co:443"
	DefaultGasPrices              = "0.1udvpn"
	DefaultChainID                = "sentinelhub-2"
	DefaultGas                    = 200000
	DefaultGasAdjustment          = 1.05
	DefaultHandshakePeers         = 8
	DefaultIntervalSetSessions    = "2m0s"
	DefaultIntervalUpdateSessions = "1h55m0s"
	DefaultIntervalUpdateStatus   = "55m0s"
	DefaultMoniker                = "My dVNP Node"
	DefaultPrice                  = "0.1tsent"
)

type Configurations struct {
	ConfPath string
	CertPath string
	DVPN     dVPNConfig
}

var Config Configurations

func LoadConfig() error {
	Config.ConfPath = appconf.Paths.DVPNConfigFullPath()
	confBytes, readErr := ioutil.ReadFile(Config.ConfPath)

	// if config does not exist, create an empty config
	if readErr != nil {
		if err := os.MkdirAll(appconf.Paths.SentinelPath(), os.ModePerm); err != nil {
			return err
		}

		Config.DVPN.Keyring.Backend = BackEndTypeTest
		Config.DVPN.Chain.Gas = DefaultGas
		Config.DVPN.Chain.GasAdjustment = DefaultGasAdjustment
		Config.DVPN.Chain.RPCAddress = DefaultRPCAddress
		Config.DVPN.Chain.GasPrices = DefaultGasPrices
		Config.DVPN.Chain.ID = DefaultChainID
		Config.DVPN.Handshake.Enable = true
		Config.DVPN.Handshake.Peers = DefaultHandshakePeers
		Config.DVPN.Node.IntervalSetSessions = "2m0s"
		Config.DVPN.Node.IntervalSetSessions = DefaultIntervalSetSessions
		Config.DVPN.Node.IntervalUpdateSessions = DefaultIntervalUpdateSessions
		Config.DVPN.Node.IntervalUpdateStatus = DefaultIntervalUpdateStatus
		Config.DVPN.Node.Moniker = DefaultMoniker
		Config.DVPN.Node.Price = DefaultPrice

		_, err := Config.PostConfig(Config.DVPN)

		return err
	}

	if unmarshalErr := toml.Unmarshal(confBytes, &Config.DVPN); unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

func GetConfigs() (config []byte, err error) {
	wgConfigPath := appconf.Paths.WireGuardConfigFullPath()
	_, readErr := ioutil.ReadFile(wgConfigPath)

	if readErr != nil {
		return initWireguardConfig()
	}

	tlsCertPath := appconf.Paths.CertificateFullPath()
	_, readErr = ioutil.ReadFile(tlsCertPath)

	if readErr != nil {
		return generateCertificate()
	}

	return config, err
}

func (c *Configurations) PostConfig(config dVPNConfig) (resp []byte, err error) {
	configPath := appconf.Paths.DVPNConfigFullPath()
	configBytes, err := toml.Marshal(config)

	if err != nil {
		return resp, err
	}

	if err = ioutil.WriteFile(configPath, configBytes, 0644); err != nil {
		return resp, err
	}

	resp, err = json.Marshal(config)

	if err != nil {
		return resp, err
	}

	c.DVPN = config

	return resp, err
}

func initWireguardConfig() (config []byte, err error) {
	cmd := exec.Command(node.DVPNNodeExec, node.DVPNNodeWireguard, node.DVPNNodeConfig, node.DVPNNodeInit)

	err = cmd.Run()

	if err != nil {
		return config, err
	}

	return GetConfigs()
}

func generateCertificate() (config []byte, err error) {
	cmd := exec.Command(node.BinSH, os.Getenv("HOME")+node.SHGenerateCertPath)

	err = cmd.Run()

	if err != nil {
		return config, err
	}

	return GetConfigs()
}
