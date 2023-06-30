package server

import (
	"go-metricscol/internal/server/router"
	"net/http"
)

func New(addr string) *http.Server {
	r := router.New()

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
