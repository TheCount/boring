package wallet

import (
	"fmt"
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

	// openWallets is the set of wallets currently open.
	openWallets map[string]*Wallet
}

// NewManager creates a new wallet manager from the specified
// configuration.
func NewManager(cfg *config.WalletConfig) (*Manager, error) {
	result := &Manager{
		cfg:         cfg,
		openWallets: make(map[string]*Wallet),
	}

	// Open all known wallets
	success := false
	for _, name := range cfg.Names {
		wallet, err := OpenWallet(cfg.WalletsDir, name)
		if err != nil {
			return nil, fmt.Errorf("Unable to open wallet '%s': %s", name, err)
		}
		defer func(wallet *Wallet) {
			if !success {
				wallet.Close()
			}
		}(wallet)
		result.openWallets[name] = wallet
	}

	success = true
	return result, nil
}

// GetWalletNames returns a copy of the list of wallet names.
func (m *Manager) GetWalletNames() (result []string) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	result = append(result, m.cfg.Names...)
	return
}

// CreateWallet creates a new wallet with the specified name and passphrase.
func (m *Manager) CreateWallet(name, password string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if m.cfg.HasName(name) {
		return fmt.Errorf("A wallet named '%s' already exists", name)
	}
	m.cfg.AddName(name)
	wallet, err := NewWallet(m.cfg.WalletsDir, name, password)
	if err != nil {
		m.cfg.RemoveName(name)
		return fmt.Errorf("Unable to create new wallet '%s': %s", name, err)
	}
	m.openWallets[name] = wallet
	if err := m.cfg.SaveNames(); err != nil {
		return fmt.Errorf(
			"Wallet '%s' successfully created, but could not save wallet name: %s",
			name, err,
		)
	}
	return nil
}

// CloseAllWallets closes all wallets opened by this wallet manager.
func (m *Manager) CloseAllWallets() (err error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	for _, wallet := range m.openWallets {
		if err2 := wallet.Close(); err2 != nil {
			err = err2
		}
	}
	m.openWallets = make(map[string]*Wallet)
	if err != nil {
		err = fmt.Errorf("Error closing wallets, most recent error was: %s", err)
	}
	return
}
