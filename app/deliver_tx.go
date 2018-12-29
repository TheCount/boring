package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// DeliverTx checks whether the specified transaction should be added
// to the current block.
func (a *App) DeliverTx(tx []byte) abcitypes.ResponseDeliverTx {
	// reject everything for now
	return abcitypes.ResponseDeliverTx{
		Code: TxInvalidType,
		Log:  "Transaction currently not supported",
	}
}
