package http

import (
	"github.com/go-chi/chi"

	"go-metricscol/internal/server/metrics"
	"go-metricscol/internal/server/middleware"
)

func MapMetricsRoutes(r *chi.Mux, h metrics.HttpHandlers, mw *middleware.Manager) {
	r.Get("/value/{type}/{name}", h.Find)
	r.Post("/value/", h.FindJSON)
	r.Post("/update/{type}/{name}/{value}", mw.DiskSaverHttpMiddleware(h.Update))
	r.Post("/update/", mw.ValidateHashHandler(mw.DiskSaverHttpMiddleware(h.UpdateJSON)))
	r.Post("/updates/", mw.ValidateHashHandler(mw.DiskSaverHttpMiddleware(h.Updates)))

	r.HandleFunc("/", h.GetAll)
}
