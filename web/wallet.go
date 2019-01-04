package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var walletTemplate = template.Must(template.New("wallet").Parse(
	`<!DOCTYPE html>
  <html>
    <head>
      <title>Boring — Wallet {{.Name}}</title>
    </head>
    <body>
      <a href="/wallets">Back</a>
      <h1>{{.Name}}: Wallet Status</h1>
      <table border>
        <tr>
          <td>Locked?</td>
          <td>{{.Locked}}</td>
        </tr>
      </table>
      <h1>Menu</h1>
      <ul>
        <li>
          {{if .Locked -}}
            <a href="/wallets/unlock/{{.Name}}">Unlock</a>
          {{- else -}}
            <a href="/wallets/lock/{{.Name}}">Lock</a>
          {{- end}}
        </li>
      </ul>
    </body>
  </html>`,
))

// walletData encapsulates the data necessary to render the wallet template.
type walletData struct {
	// Name is the wallet name.
	Name string

	// Locked indicates whether the wallet is currently locked.
	Locked bool
}

// walletHandler renders a wallet page.
type walletHandler struct {
	// web is a back pointer to the web frontend
	web *Web
}

// ServeHTTP serves a wallet page.
func (h *walletHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.web.SetDefaultHeaders(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Extract data
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	name := parts[3]
	if !h.web.WalletManager.HasWallet(name) {
		w.WriteHeader(http.StatusGone)
		return
	}
	data := walletData{
		Name:   name,
		Locked: h.web.WalletManager.IsLocked(name),
	}
	// Execute template
	if err := walletTemplate.Execute(w, data); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
	}
}

var walletUnlockTemplate = template.Must(template.New("walletUnlock").Parse(
	`<!DOCTYPE html>
  <html>
    <head>
      <title>Boring — Unlock wallet {{.}}</title>
    </head>
    <body>
      <a href="/wallets/wallet/{{.}}">Back</a>
      <h1>Unlock wallet</h1>
      <form action="/wallets/unlock/" method="post">
        <label for="pw">Wallet passphrase:</label>
        <br />
        <input type="password" id="pw" name="Password" />
        <br />
        <input type="hidden" id="name" name="Name" value="{{.}}" />
        <input type="submit" value="Unlock" />
      </form>
    </body>
  </html>`,
))

// walletUnlockRPC encapsulates data sent with a POST request
// for unlocking a wallet.
type walletUnlockRPC struct {
	// Name is the name of the wallet.
	Name string

	// Password is the wallet passphrase.
	Password string
}

// Validate validates this RPC.
func (rpc *walletUnlockRPC) Validate() error {
	if rpc.Name == "" {
		return errors.New("Wallet name must be non-empty")
	}
	return nil
}

// walletUnlockHandler handles the unlocking of wallets.
type walletUnlockHandler struct {
	// web is a back pointer to the web frontend.
	web *Web
}

// ServeHTTP serves the wallet unlock page.
func (h *walletUnlockHandler) ServeHTTP(
	w http.ResponseWriter, r *http.Request,
) {
	h.web.SetDefaultHeaders(w)
	switch r.Method {
	case http.MethodGet:
		h.writeForm(w, r)
	case http.MethodPost:
		h.unlockWallet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// writeForm writes the wallet unlock form.
func (h *walletUnlockHandler) writeForm(
	w http.ResponseWriter, r *http.Request,
) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	name := parts[3]
	if !h.web.WalletManager.HasWallet(name) {
		w.WriteHeader(http.StatusGone)
		return
	}
	if err := walletUnlockTemplate.Execute(w, name); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
	}
}

var walletUnlockedTemplate = template.Must(template.New("walletUnlocked").Parse(
	`<!DOCTYPE html>
  <html>
    <head>
      <title>Boring — Wallet {{.}} successfully unlocked</title>
    </head>
    <body>
      <a href="/wallets/wallet/{{.}}">Back</a>
      <h1>Wallet unlocked</h1>
      Wallet
      <a href="/wallets/wallet/{{.}}">{{.}}</a>
      has been successfully unlocked.
    </body>
  </html>`,
))

// unlockWallet unlocks a wallet.
func (h *walletUnlockHandler) unlockWallet(
	w http.ResponseWriter, r *http.Request,
) {
	var rpc walletUnlockRPC
	format, err := decodeRPC(r, &rpc)
	if err != nil {
		h.errorPage(w, http.StatusBadRequest, format, err)
		return
	}
	if err := h.web.WalletManager.UnlockWallet(
		rpc.Name, rpc.Password,
	); err != nil {
		h.errorPage(w, http.StatusUnauthorized, format, err)
		return
	}
	if err := walletUnlockedTemplate.Execute(w, rpc.Name); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
	}
}

var walletUnlockErrorTemplate = template.Must(template.New(
	"walletUnlockError",
).Parse(
	`<!DOCTYPE html>
	<html>
		<head>
			<title>Boring — Error unlocking</title>
		</head>
		<body>
			<a href="/wallets">Back</a>
			<h1>Error unlocking wallet</h1>
			{{.}}
		</body>
	</html>`,
))

// errorPage renders the error page for wallet unlocking.
func (h *walletUnlockHandler) errorPage(
	w http.ResponseWriter, status int, format rpcFormat, rpcErr error,
) {
	var contentType string
	switch format {
	case rpcFormatForm:
		contentType = "text/html; charset=utf-8"
	default:
		contentType = "application/json"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	switch format {
	case rpcFormatForm:
		if err := walletUnlockErrorTemplate.Execute(w, rpcErr.Error()); err != nil {
			w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
		}
	default:
		w.Write([]byte(fmt.Sprintf(`{"Error":"%s"}`, rpcErr)))
	}
}

var walletLockedTemplate = template.Must(template.New("walletLocked").Parse(
	`<!DOCTYPE html>
	<html>
		<head>
			<title>Boring — Wallet {{.}} locked</title>
		</head>
		<body>
			<a href="/wallets/wallet/{{.}}">Back</a>
			<h1>Wallet locked</h1>
			Wallet
			<a href="/wallets/wallet/{{.}}">{{.}}</a>
			has been successfully locked.
		</body>
	</html>`,
))

var walletLockErrTemplate = template.Must(template.New("walletLockErr").Parse(
	`<!DOCTYPE html>
	<html>
		<head>
			<title>Boring — Error locking wallet</title>
		</head>
		<body>
			<a href="/wallets">Back</a>
			<h1>Error unlocking wallet</h1>
			{{.}}
		</body>
	</html>`,
))

// walletLockHandler handles the locking of a wallet.
type walletLockHandler struct {
	// web is a back pointer to the web frontend.
	web *Web
}

// ServeHTTP serves the wallet lock response page.
func (h *walletLockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.web.SetDefaultHeaders(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// get name
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	name := parts[3]
	// Lock wallet
	if err := h.web.WalletManager.LockWallet(name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err = walletLockErrTemplate.Execute(w, err.Error()); err != nil {
			w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
		}
		return
	}
	if err := walletLockedTemplate.Execute(w, name); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
	}
}
