package app

// CheckTx/DeliverTx response codes.
const (
	// TxSuccess means the transaction was checked/delivered successfully.
	TxSuccess = iota

	// TxInvalidType means the transaction type is invalid.
	TxInvalidType
)
