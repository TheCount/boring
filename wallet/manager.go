package wallet

import (
	"sync"

	"github.com/TheCount/boring/config"
)

// Manager manages the set of wallets.
// There should be only one instance.
type Manager struct {
	// mtx serialises access to the wallets.
	mtx sync.Mutex

	// cfg is a pointer to the wallet configuration.
	cfg *config.WalletConfig
}

// NewManager creates a new wallet manager from the specified
// configuration.
func NewManager(cfg *config.WalletConfig) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

// GetWalletNames returns a copy of the list of wallet names.
func (m *Manager) GetWalletNames() (result []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	result = append(result, m.cfg.Names...)
	return
}
