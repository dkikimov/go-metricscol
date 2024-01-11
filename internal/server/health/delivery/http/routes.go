package http

import (
	"github.com/go-chi/chi"

	"go-metricscol/internal/server/health"
)

func MapHealthRoutes(r *chi.Mux, h health.HTTPHandlers) {
	r.Get("/ping", h.Ping)
}
