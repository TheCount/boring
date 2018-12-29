// Package app contains the implementation of Boring as a Tendermint
// application.
package app

import (
	"fmt"

	"github.com/TheCount/boring/config"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// App contains the tendermint application data.
type App struct {
	abcitypes.BaseApplication

	// State is the current application state.
	State *State
}

// NewApp creates a new application.
func NewApp(cfg *config.AppConfig) (*App, error) {
	state, err := LoadState(cfg.DBDir)
	if err != nil {
		return nil, fmt.Errorf("Unable to load application state: %s", err)
	}
	return &App{
		BaseApplication: *abcitypes.NewBaseApplication(),
		State:           state,
	}, nil
}
