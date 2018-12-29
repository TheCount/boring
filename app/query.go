package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// Query queries the app.
func (a *App) Query(req abcitypes.RequestQuery) abcitypes.ResponseQuery {
	// simply forward query to QueryState
	return a.State.QueryState.Query(req)
}
