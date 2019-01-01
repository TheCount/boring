// Package web provides the webserver for information and wallet access.
package web

import (
	"github.com/TheCount/boring/wallet"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

// Web encapsulates the web frontend.
type Web struct {
	// Client is a tendermint RPC client.
	Client tmclient.Client

	// WalletManager is the wallet manager.
	WalletManager *wallet.Manager
}

// NewWeb creates a new web frontend.
func NewWeb(client tmclient.Client, walletManager *wallet.Manager) *Web {
	return &Web{
		Client:        client,
		WalletManager: walletManager,
	}
}
