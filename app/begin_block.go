package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// BeginBlock signals the start of the delivery of a new block.
func (a *App) BeginBlock(
	req abcitypes.RequestBeginBlock,
) abcitypes.ResponseBeginBlock {
	a.DeliveredBlockHeight = req.Header.Height
	a.DeliveredBlockTime = req.Header.Time
	return abcitypes.ResponseBeginBlock{}
}
