package web

import (
	"net/http"
)

// SetDefaultHeaders sets default headers on a response writer.
func (w *Web) SetDefaultHeaders(writer http.ResponseWriter) {
	// Deactivate caching.
	writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
}
