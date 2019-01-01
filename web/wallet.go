package web

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

var walletsTemplate = template.Must(template.New("wallets").Parse(
	`<!DOCTYPE html>
  <html>
    <head>
      <title>Boring — Wallets</title>
    </head>
    <body>
      <a href="/">Back</a>
      <h1>Available Wallets</h1>
      {{with .WalletNames -}}
        <ul>
          {{range . -}}
            <li>
              <a href="/wallets/wallet/{{.}}">{{.}}</a>
            </li>
          {{- end}}
        </ul>
      {{- else -}}
        No wallets found
      {{- end}}
      <h1>Menu</h1>
      <ul>
        <li>
          <a href="/wallets/new">Create a new wallet</a>
        </li>
      </ul>
    </body>
  </html>`,
))

// walletsData encapsulates the data necessary to render the wallets template.
type walletsData struct {
	// WalletNames is the list of names of available wallets.
	WalletNames []string
}

// walletsHandler renders the wallets page.
type walletsHandler struct {
	// web is a back pointer to the web frontend.
	web *Web
}

// ServeHTTP serves the wallets page.
func (h *walletsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.web.SetDefaultHeaders(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data := walletsData{h.web.WalletManager.GetWalletNames()}
	if err := walletsTemplate.Execute(w, data); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
	}
}

var walletsNewPage = []byte(
	`<!DOCTYPE html>
  <html>
    <head>
      <title>Boring — New wallet</title>
    </head>
    <body>
      <a href="/wallets">Back</a>
      <h1>Create a new wallet</h1>
      <form action="/wallets/new" method="post">
        <label for="name">Wallet name:</label>
        <br />
        <input type="text" id="name" name="Name" />
        <br />
        <label for="pw">Wallet passphrase:</label>
        <br />
        <input type="password" id="pw" name="Password" />
        <br />
        <label for="pwConfirm">Confirm wallet passphrase:</label>
        <br />
        <input type="password" id="pwConfirm" name="PasswordConfirm" />
        <br />
        <input type="submit" value="Create" />
      </form>
    </body>
  </html>`,
)

// walletsNewRPC encapsulates the data sent with a POST request
// for creating a new wallet.
type walletsNewRPC struct {
	// Name is the name of the new wallet.
	Name string

	// Password is the wallet passphase.
	Password string

	// PasswordConfirm is the confirmation of the wallet passphrase.
	PasswordConfirm string
}

// Validate validates this RPC.
func (rpc *walletsNewRPC) Validate() error {
	if rpc.Name == "" {
		return errors.New("Wallet name must be non-empty")
	}
	if rpc.Password != rpc.PasswordConfirm {
		return errors.New("Passphrases do not match")
	}
	return nil
}

// walletsNewHandler handles the creation of new wallets
type walletsNewHandler struct {
	// web is a back pointer to the web frontend.
	web *Web
}

// ServeHTTP serves the new wallet creation page.
func (h *walletsNewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.web.SetDefaultHeaders(w)
	switch r.Method {
	case http.MethodGet:
		w.Write(walletsNewPage)
	case http.MethodPost:
		h.createWallet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// createWallet creates the new wallet.
func (h *walletsNewHandler) createWallet(
	w http.ResponseWriter, r *http.Request,
) {
	var rpc walletsNewRPC
	format, err := decodeRPC(r, &rpc)
	if err != nil {
		h.errorPage(w, http.StatusBadRequest, format, err)
		return
	}
	w.WriteHeader(http.StatusNotImplemented) // FIXME
}

var walletsNewErrorTemplate = template.Must(template.New(
	"templateNewError",
).Parse(
	`<!DOCTYPE html>
	<html>
		<head>
			<title>Boring — Error creating wallet</title>
		</head>
		<body>
			<a href="/wallets/new">Back</a>
			<h1>Error creating wallet</h1>
			{{.}}
		</body>
	</html>`,
))

// errorPage renders the error page for wallet creation.
func (h *walletsNewHandler) errorPage(
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
		if err := walletsNewErrorTemplate.Execute(w, rpcErr.Error()); err != nil {
			w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
		}
	default:
		w.Write([]byte(fmt.Sprintf(`{"Error":"%s"}`, rpcErr)))
	}
}
