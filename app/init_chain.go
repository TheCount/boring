package app

import (
	"github.com/TheCount/boring"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// InitChain initialises the blockchain of this application.
func (a *App) InitChain(
	req abcitypes.RequestInitChain,
) abcitypes.ResponseInitChain {
	if req.Time != boring.GenesisTime {
		panic("Bad genesis time")
	}
	if req.ChainId != boring.ChainIDTestnet {
		panic("Bad chain ID")
	}
	if len(req.Validators) > 0 {
		panic("Validators proposed by node")
	}
	if len(req.AppStateBytes) > 0 {
		panic("Nonzero application state proposed by node")
	}
	return abcitypes.ResponseInitChain{
		ConsensusParams: &abcitypes.ConsensusParams{
			BlockSize: &abcitypes.BlockSizeParams{
				MaxBytes: boring.MaxBlockBytes,
				MaxGas:   -1,
			},
			Evidence: &abcitypes.EvidenceParams{
				MaxAge: boring.MaxEvidenceAge,
			},
			Validator: &abcitypes.ValidatorParams{
				PubKeyTypes: []string{
					tmtypes.ABCIPubKeyTypeEd25519,
				},
			},
		},
		Validators: []abcitypes.ValidatorUpdate{
			abcitypes.ValidatorUpdate{
				PubKey: abcitypes.PubKey{
					Type: tmtypes.ABCIPubKeyTypeEd25519,
					Data: boring.GenesisValidator[:],
				},
				Power: 1,
			},
		},
	}
}
