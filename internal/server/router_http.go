package server

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/postgres"
	"go-metricscol/internal/server/metrics/delivery/http"
	middleware2 "go-metricscol/internal/server/middleware"

	"go-metricscol/internal/server/metrics/middleware"
)

type HttpRouter struct {
	chi.Router
}

func newHttpRouter(storage repository.Repository, postgres *postgres.DB, config *Config) HttpRouter {
	processors := http.NewHandlers(
		storage,
		postgres,
		http.NewConfig(
			config.HashKey,
			config.CryptoKey,
		),
	)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "text/plain", "text/xml"))
	r.Use(chiMiddleware.Logger)
	r.Use(middleware2.DecompressHandler)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(s.TrustedSubnetHandler)

	r.Get("/value/{type}/{name}", processors.Find)
	r.Post("/value/", processors.FindJSON)

	r.Post("/update/{type}/{name}/{value}", Conveyor(config, processors.Update, middleware.diskSaverHandler))
	r.Post("/update/", Conveyor(config, processors.UpdateJSON, middleware.diskSaverHandler, middleware2.ValidateHashHandler))
	r.Post("/updates/", Conveyor(config, processors.Updates, middleware.diskSaverHandler, middleware2.ValidateHashesHandler))

	r.Get("/ping", processors.Ping)

	r.HandleFunc("/", processors.GetAll)
	return HttpRouter{r}
}
