package server

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	helathUseCase "go-metricscol/internal/server/health/usecase"
	metricsUseCase "go-metricscol/internal/server/metrics/usecase"

	healthHttp "go-metricscol/internal/server/health/delivery/http"
	metricsHttp "go-metricscol/internal/server/metrics/delivery/http"

	"go-metricscol/internal/server/middleware"
)

func (s Server) MapHandlers(r *chi.Mux) error {
	metricsUC := metricsUseCase.NewMetricsUC(s.Repo, s.Postgres, s.Config)
	healthUC := helathUseCase.NewHealthUC(s.Postgres)

	mw := middleware.NewManager(metricsUC, healthUC, s.Config, s.Repo)

	r.Use(chiMiddleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "text/plain", "text/xml"))
	r.Use(chiMiddleware.Logger)
	r.Use(mw.DecompressHandler)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(mw.TrustedSubnetHandler)

	metricsHttp.NewMetricsHandlers(metricsUC, s.Config)
	healthHttp.NewHealthHandlers(healthUC)

	return nil
}
