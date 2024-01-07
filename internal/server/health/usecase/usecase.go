package usecase

import (
	"context"

	"go-metricscol/internal/repository"
)

type HealthUC struct {
	repository repository.Repository
}

func NewHealthUC(repository repository.Repository) *HealthUC {
	return &HealthUC{repository: repository}
}

func (h HealthUC) Ping(ctx context.Context) error {
	return h.repository.Ping(ctx)
}
