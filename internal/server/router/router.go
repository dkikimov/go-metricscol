package router

import (
	"github.com/go-chi/chi"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/server/handlers"
)

func New() chi.Router {
	processors := handlers.Processors{
		Storage: memory.NewMemStorage(),
	}

	r := chi.NewRouter()

	r.HandleFunc("/update/{type}/{name}/{value}", processors.Update)
	r.HandleFunc("/", processors.NotFound)
	return r
}
