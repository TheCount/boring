package boring

import (
	"time"
)

const (
	// ChainID is the name of the blockchain.
	ChainID = "boring"

	// ChainIDTestnet is the name of the blockchain (on testnet).
	ChainIDTestnet = "boring-testnet"

	// FeeDecreaseThreshold is the total transactions size threshold
	// in bytes below which the transaction fee per byte decreases for
	// the next block.
	FeeDecreaseThreshold = 100 * 1024

	// FeeIncreaseThreshold is the total transactions size threshold
	// in bytes above which no more transactions are accepted for the current
	// block and the transaction fee per byte increases for the next block.
	FeeIncreaseThreshold = 1024 * 1024

	// MaxBlockBytes is the maximum number of bytes in a block.
	MaxBlockBytes = 2 * 1024 * 1024

	// MaxEvidenceAge is the maximum number of blocks between the current block
	// and a block in which validator malfeance occurred.
	MaxEvidenceAge = 600000

	// MaxTransactionBytes is the maximum size of one transaction in bytes.
	MaxTransactionBytes = 10 * 1024
)

var (
	// GenesisTime is the blockchain genesis time.
	GenesisTime = time.Date(2018, 12, 25, 13, 0, 0, 0, time.UTC)
)
