package middleware

import (
	"compress/gzip"
	"go-metricscol/internal/server/apierror"
	"net/http"
)

func DecompressHandler(next http.Handler) http.Handler {
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
