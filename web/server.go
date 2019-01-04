package web

import (
	"net/http"
)

// Serve starts the web server and blocks until an error occurs.
func (w *Web) Serve() error {
	http.Handle("/", &mainHandler{w})
	http.Handle("/wallets/", &walletsHandler{w})
	http.Handle("/wallets/lock/", &walletLockHandler{w})
	http.Handle("/wallets/new", &walletsNewHandler{w})
	http.Handle("/wallets/unlock/", &walletUnlockHandler{w})
	http.Handle("/wallets/wallet/", &walletHandler{w})
	return http.ListenAndServe("localhost:22222", nil)
}
