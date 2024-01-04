package usecase

import (
	"context"

	"go-metricscol/internal/repository/postgres"
)

type HealthUC struct {
	Postgres *postgres.DB
}

func NewHealthUC(postgres *postgres.DB) *HealthUC {
	return &HealthUC{Postgres: postgres}
}

func (h HealthUC) Ping(ctx context.Context) error {
	return h.Postgres.Ping(ctx)
}
