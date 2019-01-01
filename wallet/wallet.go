// Package wallet does all the wallet handling.
package wallet

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/libs/db"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

// Prefixes for the various types of data within the wallet.
const (
	// prefixPassword is the key used for storing the salted and hashed wallet
	// password.
	prefixPassword byte = iota

	// prefixKey is the key used for storing the encrypted master key
	// of the wallet.
	prefixKey
)

const (
	// passwordHashCost is the cost used for bcrypt hashing of the password.
	passwordHashCost = 12

	// keyLength is the master key length in bytes.
	keyLength = 32
)

// Wallet encapsulates a wallet.
type Wallet struct {
	// store is the storage backend of this wallet.
	store db.DB
}

// NewWallet creates a new wallet with the specified name in the specified
// directory. The secret parts of the new wallet will be encrypted using a key
// derived from the specified password.
func NewWallet(dir, name, password string) (*Wallet, error) {
	// Create wallet store.
	store, err := db.NewGoLevelDB(name, dir)
	if err != nil {
		return nil, fmt.Errorf("Unable to create wallet database: %s", err)
	}
	success := false
	defer func() {
		if !success {
			store.Close()
		}
	}()
	if store.Iterator(nil, nil).Valid() {
		return nil, fmt.Errorf("Database '%s/%s' already exists", dir, name)
	}

	// Hash and salt password and encrypt master key.
	nonce := make([]byte, keyLength)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, errors.New("Unable to generate random nonce")
	}
	hashedPW, err := bcrypt.GenerateFromPassword(
		[]byte(password), passwordHashCost,
	)
	if err != nil {
		return nil, errors.New("Unable to hash wallet password")
	}
	batch := store.NewBatch()
	batch.Set([]byte{prefixPassword}, hashedPW)
	tempKey := sha3.Sum256(bytes.Join([][]byte{
		[]byte(password), hashedPW,
	}, nil))
	cipher, err := aes.NewCipher(tempKey[:])
	if err != nil {
		return nil, fmt.Errorf("Unable to create AES-256 cipher: %s", err)
	}
	encryptedKey := make([]byte, keyLength)
	cipher.Encrypt(encryptedKey, nonce)
	batch.Set([]byte{prefixKey}, encryptedKey)
	batch.WriteSync()

	success = true
	return &Wallet{
		store: store,
	}, nil
}

// OpenWallet opens the specified wallet.
func OpenWallet(dir, name string) (*Wallet, error) {
	// Open wallet store.
	store, err := db.NewGoLevelDB(name, dir)
	if err != nil {
		return nil, fmt.Errorf("Unable to open wallet database: %s", err)
	}
	if !store.Iterator(nil, nil).Valid() {
		store.Close()
		return nil, fmt.Errorf("Database '%s/%s' does not exist", dir, name)
	}
	return &Wallet{
		store: store,
	}, nil
}

// Close closes this wallet.
// Regardless of the return value, it is illegal to refer to this wallet
// after it has been closed.
func (w *Wallet) Close() error {
	w.store.Close()
	return nil
}
