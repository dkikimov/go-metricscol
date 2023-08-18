package handlers

import (
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/postgres"
)

type Handlers struct {
	Storage  repository.Repository
	Postgres *postgres.DB
	config   *Config
}

func NewHandlers(storage repository.Repository, postgres *postgres.DB, config *Config) *Handlers {
	return &Handlers{Storage: storage, Postgres: postgres, config: config}
}
