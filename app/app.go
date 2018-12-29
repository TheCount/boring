// Package app contains the implementation of Boring as a Tendermint
// application.
package app

import (
	"fmt"
	"time"

	"github.com/TheCount/boring/config"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

const (
	// Version is the application semantic version.
	Version = "0.0.0"

	// ProtocolVersion is the current protocol version of the application.
	ProtocolVersion = 0

	// StoreFormatVersion is the version of the storage format.
	StoreFormatVersion = 0
)

// App contains the tendermint application data.
type App struct {
	abcitypes.BaseApplication

	// State is the current application state.
	State *State

	// DeliveredBlockHeight is the height of the block currently being processed.
	// To be used only by the consensus connection.
	DeliveredBlockHeight int64

	// DeliveredBlockTime is the time of the block currently being processed.
	// To be used only by the consensus connection.
	DeliveredBlockTime time.Time
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
