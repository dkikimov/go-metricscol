package router

import (
	"github.com/go-chi/chi"
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
