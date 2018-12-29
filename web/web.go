// Package web provides the webserver for information and wallet access.
package web

import (
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

// Web encapsulates the web frontend.
type Web struct {
	// Client is a tendermint RPC client.
	Client tmclient.Client
}

// NewWeb creates a new web frontend.
func NewWeb(client tmclient.Client) *Web {
	return &Web{client}
}
