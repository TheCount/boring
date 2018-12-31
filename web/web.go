// Package web provides the webserver for information and wallet access.
package web

import (
	"sync"

	"github.com/TheCount/boring/config"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

// Web encapsulates the web frontend.
type Web struct {
	// mtx is a mutex to protect those members of this struct which are
	// not concurency safe.
	mtx sync.Mutex

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

// GetWalletNames obtains a copy of the current wallet names.
func (w *Web) GetWalletNames() (result []string) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	result = append(result, w.WalletConfig.Names...)
	return
}
