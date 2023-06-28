package server

import (
	"go-metricscol/internal/server/handlers"
	"net/http"
)

func GetServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/Gauge/", handlers.UpdateGauge)
	mux.HandleFunc("/update/Counter/", handlers.UpdateCounter)

	return mux
}
