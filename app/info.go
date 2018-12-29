package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// Info obtains application information.
func (a *App) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	commit := a.State.QueryState.LastCommitID()
	return abcitypes.ResponseInfo{
		Version:          Version,
		AppVersion:       ProtocolVersion,
		LastBlockHeight:  commit.Version,
		LastBlockAppHash: commit.Hash,
	}
}
