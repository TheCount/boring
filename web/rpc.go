package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// rpcFormat encodes the RPC format.
type rpcFormat uint8

// Supported RPC formats.
const (
	// rpcFormatForm denotes RPCs using WWW forms.
	rpcFormatForm rpcFormat = iota

	// rpcFormatKSOM denotes JSON-encoded RPCs.
	rpcFormatJSON
)

// rpc is an interface for generic web RPCs.
type rpc interface {
	// Validate validates the RPC.
	Validate() error
}

// decodeRPC decodes the specified rpc from the specified request.
func decodeRPC(r *http.Request, rpc rpc) (rpcFormat, error) {
	// Check whether request is WWW form encoded.
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		return rpcFormatForm, decodeForm(r, rpc)
	}
	// In all other cases, assume it's a JSON request.
	return rpcFormatJSON, decodeJSON(r, rpc)
}

// decodeForm decodes the specified WWW form request into the specified RPC.
func decodeForm(r *http.Request, rpc rpc) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	structure := make(map[string]string)
	for key, values := range r.PostForm {
		if len(values) == 0 {
			continue
		}
		if len(values) > 1 {
			return fmt.Errorf("POST key '%s' has multiple values", key)
		}
		structure[key] = values[0]
	}
	if err := mapstructure.WeakDecode(structure, rpc); err != nil {
		return fmt.Errorf("Unable to decode form: %s", err)
	}
	return rpc.Validate()
}

// decodeJSON interprets the specified request as JSON and unmarshals it
// into the specified RPC.
func decodeJSON(r *http.Request, rpc rpc) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rpc); err != nil {
		return fmt.Errorf("Error unmarshaling JSON: %s", err)
	}
	return rpc.Validate()
}
