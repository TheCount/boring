package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// EndBlock signals the end of the delivery of a new block.
func (a *App) EndBlock(
	req abcitypes.RequestEndBlock,
) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}
