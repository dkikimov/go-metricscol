package backends

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"go-metricscol/internal/config"
	"go-metricscol/internal/repository"
	healthHttp "go-metricscol/internal/server/health/delivery/http"
	healthUseCase "go-metricscol/internal/server/health/usecase"
	metricsUseCase "go-metricscol/internal/server/metrics/usecase"

	metricsHttp "go-metricscol/internal/server/metrics/delivery/http"
	"go-metricscol/internal/server/middleware"
)

type HTTP struct {
	server *http.Server
}

func NewHTTP(repo repository.Repository, config *config.ServerConfig) (*HTTP, error) {
	r := chi.NewRouter()

	metricsUC := metricsUseCase.NewMetricsUC(repo, config)
	healthUC := healthUseCase.NewHealthUC(repo)

	mw := middleware.NewManager(metricsUC, healthUC, config, repo)

	r.Use(chiMiddleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "text/plain", "text/xml"))
	r.Use(chiMiddleware.Logger)
	r.Use(mw.DecompressHandler)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(mw.HTTPTrustedSubnetHandler)

	healthHttp.NewHealthHandlers(healthUC)

	metricsHttp.MapMetricsRoutes(r, metricsHttp.NewMetricsHandlers(metricsUC, config), mw)
	healthHttp.MapHealthRoutes(r, healthHttp.NewHealthHandlers(healthUC))

	httpServer := http.Server{
		Addr:    config.Address,
		Handler: r,
	}

	return &HTTP{server: &httpServer}, nil
}

func (s HTTP) ListenAndServe() error {
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s HTTP) GracefulShutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
