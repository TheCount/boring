package web

import (
	"net/http"
)

// Serve starts the web server and blocks until an error occurs.
func (w *Web) Serve() error {
	http.Handle("/", &mainHandler{w})
	http.Handle("/wallets/", &walletsHandler{w})
	return http.ListenAndServe("localhost:22222", nil)
}
