package middleware

import (
	"bytes"
	"encoding/json"
	"go-metricscol/internal/models"
	"io"
	"net/http"
)

func ValidateHashHandler(next http.HandlerFunc, key string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(key) != 0 {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "couldn't read body", http.StatusInternalServerError)
				return
			}

			var metric models.Metric
			err = json.Unmarshal(body, &metric)
			if err != nil {
				http.Error(w, "couldn't parse json", http.StatusInternalServerError)
				return
			}

			if metric.HashValue(key) != metric.Hash {
				http.Error(w, "hash mismatch", http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		next.ServeHTTP(w, r)
	})

}
