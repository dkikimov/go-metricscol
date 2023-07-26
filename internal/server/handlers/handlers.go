package handlers

import (
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/postgres"
)

type Handlers struct {
	Storage  repository.Repository
	Postgres *postgres.DB
}
