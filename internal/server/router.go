package server

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/handlers"
	"go-metricscol/internal/server/middleware"
	"log"
	"net/http"
)

func (s Server) diskSaverHandler(next http.HandlerFunc, saveToDisk bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if saveToDisk {
			if err := s.saveToDisk(); err != nil {
				log.Printf("Couldn't save metrics to disk with error: %s", err)
			}
		}
	}
}

func (s Server) newRouter(storage repository.Repository) chi.Router {
	processors := handlers.Handlers{
		Storage:  storage,
		Postgres: s.Postgres,
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "text/plain", "text/xml"))
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.DecompressHandler)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))

	saveToDisk := s.Config.StoreInterval == 0 && len(s.Config.StoreFile) != 0

	r.Post("/update/{type}/{name}/{value}", s.diskSaverHandler(processors.Update, saveToDisk))
	r.Get("/value/{type}/{name}", processors.Get)

	r.Post("/update/", middleware.ValidateHashHandler(s.diskSaverHandler(processors.UpdateJSON, saveToDisk), s.Config.HashKey))
	r.Post("/value/", processors.GetJSON)

	r.Get("/ping", processors.Ping)

	r.HandleFunc("/", processors.GetAll)
	return r
}
