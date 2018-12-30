package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/shibukawa/configdir"
)

// WalletConfig holds the wallet configuration.
type WalletConfig struct {
	// WalletsDir is the wallets directory
	WalletsDir string

	// namesPath is the path to the names file.
	namesPath string

	// Names is a list of wallet names
	Names []string
}

// getWalletConfig obtains the wallets configuration.
func getWalletConfig() (*WalletConfig, error) {
	var result WalletConfig
	dir := configdir.New("boring", "wallet")
	basefolder := dir.QueryFolders(configdir.Global)[0]
	result.WalletsDir = filepath.Join(basefolder.Path, "wallets")
	if err := os.MkdirAll(result.WalletsDir, 0700); err != nil {
		return nil, fmt.Errorf("Unable to create wallets directory: %s", err)
	}
	result.namesPath = filepath.Join(basefolder.Path, "names.json")
	if err := result.LoadNames(); err != nil {
		return nil, err
	}
	return &result, nil
}

// LoadNames loads the wallet names from a names file.
// If the file does not exist yet, the resulting name list will be empty
// (this is not an error).
func (wc *WalletConfig) LoadNames() error {
	data, err := ioutil.ReadFile(wc.namesPath)
	if os.IsNotExist(err) {
		wc.Names = nil
		return nil
	}
	if err != nil {
		return fmt.Errorf("Unable to read wallet names file: %s", err)
	}
	if err = json.Unmarshal(data, &wc.Names); err != nil {
		return fmt.Errorf("Unable to unmarshal wallet names: %s", err)
	}
	sort.Strings(wc.Names)
	return nil
}

// SaveNames saves the wallet names in the names file.
func (wc *WalletConfig) SaveNames() error {
	data, err := json.Marshal(wc.Names)
	if err != nil {
		return fmt.Errorf("Unable to marshal wallet names: %s", err)
	}
	if err = ioutil.WriteFile(wc.namesPath, data, 0666); err != nil {
		return fmt.Errorf("Unable to write wallet names file: %s", err)
	}
	return nil
}

// AddName adds a wallet name to this configuration.
// If the name already exists, no operation is performed.
func (wc *WalletConfig) AddName(name string) {
	index := sort.SearchStrings(wc.Names, name)
	if index == len(wc.Names) || wc.Names[index] != name {
		wc.Names = append(wc.Names, "")
		copy(wc.Names[index+1:], wc.Names[index:])
		wc.Names[index] = name
	}
}

// RemoveName removes a wallet name from this configuration.
// If the name does not exist, no operation is performed.
func (wc *WalletConfig) RemoveName(name string) {
	index := sort.SearchStrings(wc.Names, name)
	if index != len(wc.Names) && wc.Names[index] == name {
		copy(wc.Names[index:], wc.Names[index+1:])
		wc.Names = wc.Names[:len(wc.Names)-1]
	}
}
