package server

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/handlers"
	"go-metricscol/internal/server/middleware"
)

func (s Server) newRouter(storage repository.Repository) chi.Router {
	processors := handlers.NewHandlers(
		storage,
		s.Postgres,
		handlers.NewConfig(
			s.Config.HashKey,
			s.Config.CryptoKey,
		),
	)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "text/plain", "text/xml"))
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.DecompressHandler)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(s.TrustedSubnetHandler)

	r.Get("/value/{type}/{name}", processors.Find)
	r.Post("/value/", processors.FindJSON)

	r.Post("/update/{type}/{name}/{value}", Conveyor(s.Config, processors.Update, s.diskSaverHandler))
	r.Post("/update/", Conveyor(s.Config, processors.UpdateJSON, s.diskSaverHandler, ValidateHashHandler))
	r.Post("/updates/", Conveyor(s.Config, processors.Updates, s.diskSaverHandler, ValidateHashesHandler))

	r.Get("/ping", processors.Ping)

	r.HandleFunc("/", processors.GetAll)
	return r
}
