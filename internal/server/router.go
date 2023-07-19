package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/handlers"
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
		Storage: storage,
	}

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "text/plain", "text/xml"))
	r.Use(middleware.Logger)
	r.Use(decompressHandler)
	r.Use(middleware.AllowContentEncoding("gzip"))

	saveToDisk := s.Config.StoreInterval == 0 && len(s.Config.StoreFile) != 0

	r.Post("/update/{type}/{name}/{value}", s.diskSaverHandler(processors.Update, saveToDisk))
	r.Get("/value/{type}/{name}", processors.Get)

	r.Post("/update/", s.diskSaverHandler(processors.UpdateJSON, saveToDisk))
	r.Post("/value/", processors.GetJSON)

	r.HandleFunc("/", processors.GetAll)
	return r
}
