package web

import (
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
	data := walletsData{h.web.GetWalletNames()}
	if err := walletsTemplate.Execute(w, data); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
	}
}
