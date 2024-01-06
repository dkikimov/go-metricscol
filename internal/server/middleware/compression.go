package middleware

import (
	"compress/gzip"
	"net/http"

	"go-metricscol/internal/server/apierror"
)

// DecompressHandler is a middleware which analyzes HTTP request header and if necessary parses Gzip request.
// http.Request body is replaced with parsed Gzip body.
func (mw *Manager) DecompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				apierror.WriteHTTP(w, err)
				return
			}

			r.Body = reader
			r.Header.Del("Content-Encoding")
			r.Header.Del("Content-Length")
		}
		next.ServeHTTP(w, r)
	})
}
