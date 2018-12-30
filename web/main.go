package web

import (
	"fmt"
	"html/template"
	"net/http"

	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var mainTemplate = template.Must(template.New("main").Parse(
	`<!DOCTYPE html>
  <html>
    <head>
      <title>Boring — Main page</title>
    </head>
    <body>
      <h1>Status</h1>
      <table border>
        <tr>
          <td>LatestBlockHeight</td>
          <td>{{.Status.SyncInfo.LatestBlockHeight}}</td>
        </tr>
        <tr>
          <td>LatestBlockTime</td>
          <td>{{.Status.SyncInfo.LatestBlockTime}}</td>
        </tr>
        <tr>
          <td>Still syncing?</td>
          <td>{{.Status.SyncInfo.CatchingUp}}</td>
        </tr>
      </table>
      <br />
      <a href="/">Refresh</a>
			<h1>Menu</h1>
			<ul>
				<li>
					<a href="/wallets">Wallets</a>
				</li>
			</ul>
    </body>
  </html>`,
))

// mainData encapsulates the data necessary to render the main template.
type mainData struct {
	// Status is the tendermint node status.
	Status *tmrpctypes.ResultStatus
}

// mainHandler renders the main page.
type mainHandler struct {
	// web is a back pointer to the web frontend.
	web *Web
}

// ServeHTTP serves the main page.
func (h *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.web.SetDefaultHeaders(w)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var data mainData
	var err error
	data.Status, err = h.web.Client.Status()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unable to obtain status: %s", err)))
		return
	}
	if err := mainTemplate.Execute(w, data); err != nil {
		w.Write([]byte(fmt.Sprintf("Template error: %s", err)))
	}
}
