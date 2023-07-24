package handlers

import (
	"github.com/jackc/pgx/v5"
	"go-metricscol/internal/repository"
)

type Handlers struct {
	Storage  repository.Repository
	Postgres *pgx.Conn
}
