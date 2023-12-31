package usecase

import (
	"context"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/postgres"
	"go-metricscol/internal/server"
)

type metricsUC struct {
	Storage  repository.Repository
	Postgres *postgres.DB
	config   *server.Config
}

func (m *metricsUC) Find(ctx context.Context, name string, mType models.MetricType) (*models.Metric, error) {
	return m.Storage.Get(ctx, name, mType)
}

func (m *metricsUC) Update(ctx context.Context, metric models.Metric) error {
	return m.Storage.Update(ctx, metric)
}

func (m *metricsUC) Updates(ctx context.Context, metrics []models.Metric) error {
	return m.Storage.Updates(ctx, metrics)
}

func (m *metricsUC) GetAll(ctx context.Context) ([]models.Metric, error) {
	return m.Storage.GetAll(ctx)
}

func NewMetricsUC(storage repository.Repository, postgres *postgres.DB, config *server.Config) *metricsUC {
	return &metricsUC{Storage: storage, Postgres: postgres, config: config}
}
