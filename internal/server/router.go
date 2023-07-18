package server

import (
	"github.com/go-chi/chi"
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
	saveToDisk := s.Config.StoreInterval == 0 && len(s.Config.StoreFile) != 0

	r.Post("/update/{type}/{name}/{value}", s.diskSaverHandler(processors.Update, saveToDisk))
	r.Get("/value/{type}/{name}", s.diskSaverHandler(processors.Get, saveToDisk))

	r.Post("/update/", s.diskSaverHandler(processors.UpdateJSON, saveToDisk))
	r.Post("/value/", s.diskSaverHandler(processors.GetJSON, saveToDisk))

	r.HandleFunc("/", s.diskSaverHandler(processors.GetAll, saveToDisk))
	return r
}
