package app

import (
	"errors"

	"github.com/TheCount/boring"
	cstore "github.com/cosmos/cosmos-sdk/store"
	ctypes "github.com/cosmos/cosmos-sdk/types"
	tmdb "github.com/tendermint/tendermint/libs/db"
)

// Store is an interface for a key-value store which can calculate commit
// hashes and which can be queried (as in ABCI Query).
type Store interface {
	cstore.Queryable
	cstore.CommitKVStore
}

// State represents the application state,
// consisting of query, check, and deliver state as suggested by the
// tendermint application development guide.
type State struct {
	// QueryState is the last committed state, used for querying.
	QueryState Store

	// CheckState is the state used by CheckTx for checking whether
	// transactions are good for the mempool.
	// Upon commit, its changes are discarded
	// and reset to the new committed state.
	CheckState cstore.CacheKVStore

	// DeliverState is the state used by DeliverTx to add transactions to the
	// current block. Upon commit, its changes are written back to the
	// QueryState.
	DeliverState cstore.CacheKVStore
}

// LoadState loads the application state from the specified database path.
func LoadState(dbPath string) (result *State, err error) {
	// Open underlying DB
	db, err := tmdb.NewGoLevelDB(boring.ChainIDTestnet, dbPath)
	if err != nil {
		return
	}

	// Make sure underlying DB is closed if anything goes wrong
	defer func() {
		p := recover()
		if p != nil || err != nil {
			db.Close()
		}
		if p != nil {
			panic(p)
		}
	}()

	// Wrap states around underlying DB
	store, err := cstore.LoadIAVLStore(
		db, cstore.CommitID{}, ctypes.PruneEverything,
	)
	if err != nil {
		return
	}
	queryState, ok := store.(Store)
	if !ok {
		err = errors.New("IAVL store does not fulfil requirements")
		return
	}
	result = &State{
		QueryState:   queryState,
		CheckState:   cstore.NewCacheKVStore(queryState),
		DeliverState: cstore.NewCacheKVStore(queryState),
	}
	return
}
