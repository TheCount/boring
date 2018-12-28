package boring

import (
	"time"

	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
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
	// GenesisValidator is the public key of the genesis validator
	GenesisValidator = tmed25519.PubKeyEd25519{
		0x93, 0x94, 0x3D, 0xE4, 0xBB, 0xE0, 0x3E, 0x0E,
		0x9B, 0xA4, 0x88, 0x7A, 0x3C, 0x54, 0x83, 0xA5,
		0x40, 0xBD, 0xF5, 0x45, 0x4A, 0x7D, 0x50, 0x5F,
		0x0D, 0x01, 0xA2, 0xF2, 0x9C, 0xD9, 0xAE, 0x5B,
	}

	// GenesisTime is the blockchain genesis time.
	GenesisTime = time.Date(2018, 12, 25, 13, 0, 0, 0, time.UTC)
)
