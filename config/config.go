package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/shibukawa/configdir"
	tmcfg "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Config holds the application configuration
type Config struct {
	// TMConfig is the tendermint node configuration.
	TMConfig *tmcfg.Config
}

// ensureTMGenesis creates the tendermint genesis configuration
// in the specified folder if it does
// not exist yet.
func ensureTMGenesis(folder *configdir.Config) error {
	genesisConfigFile := filepath.Join("config", "genesis.json")
	if folder.Exists(genesisConfigFile) {
		return nil
	}
	if err := folder.CreateParentDir(genesisConfigFile); err != nil {
		return fmt.Errorf("Unable to create genesis config directory: %s", err)
	}
	genesis := tmtypes.GenesisDoc{
		GenesisTime: time.Date(2018, 12, 25, 13, 0, 0, 0, time.UTC),
		ChainID:     "boring-testnet",
	}
	jsonBytes, err := json.Marshal(&genesis)
	if err != nil {
		return fmt.Errorf("Unable to marshal genesis config: %s", err)
	}
	if err := folder.WriteFile(genesisConfigFile, jsonBytes); err != nil {
		return fmt.Errorf("Unable to write genesis config: %s", err)
	}
	return nil
}

// getTMRootDir obtains the tendermint root directory.
func getTMRootDir() (string, error) {
	dir := configdir.New("boring", "tendermint")
	folder := dir.QueryFolders(configdir.Global)[0]
	if err := ensureTMGenesis(folder); err != nil {
		return "", err
	}
	return folder.Path, nil
}

// GetConfig obtains the configuration
func GetConfig() (*Config, error) {
	tmconfig := tmcfg.DefaultConfig()
	rootdir, err := getTMRootDir()
	if err != nil {
		return nil, err
	}
	return &Config{tmconfig.SetRoot(rootdir)}, nil
}
