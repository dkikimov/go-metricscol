package server

import (
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/server/handlers"
	"net/http"
)

func Get(addr string) *http.Server {
	processors := handlers.Processors{
		Storage: memory.NewMemStorage(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", processors.Update)
	mux.HandleFunc("/", processors.NotFound)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
