package app

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// SetOption response codes.
const (
	// SetOptionSuccess is returned if the option could be successfully set.
	SetOptionSuccess = iota

	// SetOptionInvalidKey is returned if the option key is not supported.
	SetOptionInvalidKey
)

// SetOption sets application options.
// There are currently no options,
// so a call to this method always results in failure.
func (a *App) SetOption(
	req abcitypes.RequestSetOption,
) abcitypes.ResponseSetOption {
	// No option key is currently supported.
	return abcitypes.ResponseSetOption{
		Code: SetOptionInvalidKey,
		Log:  "No options are currently supported",
	}
}
