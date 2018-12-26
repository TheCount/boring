// Package app contains the implementation of Boring as a Tendermint
// application.
package app

import (
	"github.com/TheCount/boring/config"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// App contains the tendermint application data.
type App struct {
	abcitypes.BaseApplication

	// DBDir is the App database directory path.
	DBDir string
}

// NewApp creates a new application.
func NewApp(cfg *config.AppConfig) (*App, error) {
	return &App{
		BaseApplication: *abcitypes.NewBaseApplication(),
		DBDir:           cfg.DBDir,
	}, nil
}
