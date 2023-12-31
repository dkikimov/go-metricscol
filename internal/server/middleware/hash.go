package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server"
)

// ValidateHashHandler is a middleware which gets models.Metric from request, calculates hash and compares it with given.
// If the hashes do not match, http.Error is called with code 400.
func ValidateHashHandler(next http.HandlerFunc, config *server.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashKey := config.HashKey
		if len(hashKey) != 0 {
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

			if metric.HashValue(hashKey) != metric.Hash {
				http.Error(w, "hash mismatch", http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		next.ServeHTTP(w, r)
	}
}

// ValidateHashesHandler is a middleware which gets []models.Metric from request, calculates hashes and compares them with given.
// If at least one of the hashes do not match, http.Error is called with code 400.
func ValidateHashesHandler(next http.HandlerFunc, config *server.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashKey := config.HashKey
		if len(hashKey) != 0 {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "couldn't read body", http.StatusInternalServerError)
				return
			}

			var metrics []models.Metric
			err = json.Unmarshal(body, &metrics)
			if err != nil {
				http.Error(w, "couldn't parse json", http.StatusInternalServerError)
				return
			}

			for _, metric := range metrics {
				if metric.HashValue(hashKey) != metric.Hash {
					http.Error(w, "hash mismatch", http.StatusBadRequest)
					return
				}
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		next.ServeHTTP(w, r)
	}
}
