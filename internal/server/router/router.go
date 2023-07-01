package router

import (
	"github.com/go-chi/chi"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/server/handlers"
)

func New() chi.Router {
	h := handlers.Handlers{
		Storage: memory.NewMemStorage(),
	}

	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", h.Update)
	r.Get("/value/{type}/{name}", h.Get)
	r.HandleFunc("/", h.GetAll)
	return r
}

func NewWithStorage(storage repository.Repository) chi.Router {
	processors := handlers.Handlers{
		Storage: storage,
	}

	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", processors.Update)
	r.Get("/value/{type}/{name}", processors.Get)
	r.HandleFunc("/", processors.GetAll)
	return r
}
