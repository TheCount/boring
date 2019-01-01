package web

import (
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
