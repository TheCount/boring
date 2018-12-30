// Package web provides the webserver for information and wallet access.
package web

import (
	"github.com/TheCount/boring/config"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

// Web encapsulates the web frontend.
type Web struct {
	// Client is a tendermint RPC client.
	Client tmclient.Client

	// WalletConfig is the wallet configuration.
	WalletConfig *config.WalletConfig
}

// NewWeb creates a new web frontend.
func NewWeb(client tmclient.Client, walletConfig *config.WalletConfig) *Web {
	return &Web{
		Client:       client,
		WalletConfig: walletConfig,
	}
}
