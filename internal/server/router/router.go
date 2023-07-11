package router

import (
	"github.com/go-chi/chi"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/server/handlers"
)

func New() chi.Router {
	return NewWithStorage(memory.NewMemStorage())
}

func NewWithStorage(storage repository.Repository) chi.Router {
	processors := handlers.Handlers{
		Storage: storage,
	}

	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", processors.Update)
	r.Get("/value/{type}/{name}", processors.Get)

	r.Post("/update/", processors.UpdateJSON)
	r.Post("/value/", processors.GetJSON)

	r.HandleFunc("/", processors.GetAll)
	return r
}
