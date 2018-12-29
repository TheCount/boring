package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// Commit commits the current block.
func (a *App) Commit() abcitypes.ResponseCommit {
	a.CheckMinBlockHeight = a.DeliveredBlockHeight + 1
	a.CheckBeforeBlockTime = a.DeliveredBlockTime
	commit := a.State.Commit()
	return abcitypes.ResponseCommit{
		Data: commit.Hash,
	}
}
