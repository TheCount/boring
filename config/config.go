package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/TheCount/boring"
	"github.com/shibukawa/configdir"
	tmcfg "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"
)

// Config holds the application configuration
type Config struct {
	// TMConfig is the tendermint node configuration.
	TMConfig *tmcfg.Config

	// AppConfig is the tendermint app configuration.
	AppConfig *AppConfig
}

// AppConfig holds the tendermint application configuration.
type AppConfig struct {
	// DBDir is the database directory path.
	DBDir string
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
		GenesisTime: boring.GenesisTime,
		ChainID:     boring.ChainIDTestnet,
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

// getAppDir obtains the tendermint app database directory.
func getAppDir() (string, error) {
	dir := configdir.New("boring", "app")
	folder := dir.QueryFolders(configdir.Global)[0]
	if err := folder.MkdirAll(); err != nil {
		return "", fmt.Errorf("Unable to create app directory: %s", err)
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
	appdir, err := getAppDir()
	if err != nil {
		return nil, err
	}
	return &Config{
		TMConfig:  tmconfig.SetRoot(rootdir),
		AppConfig: &AppConfig{appdir},
	}, nil
}
