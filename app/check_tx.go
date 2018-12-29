package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// CheckTx checks whether the specified transaction should be allowed to
// enter the mempool.
func (a *App) CheckTx(tx []byte) abcitypes.ResponseCheckTx {
	// reject everything for now
	return abcitypes.ResponseCheckTx{
		Code: TxInvalidType,
		Log:  "Transaction currently not supported",
	}
}
