package usecase

import (
	"context"

	"go-metricscol/internal/config"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
)

type MetricsUC struct {
	Storage repository.Repository
	config  *config.ServerConfig
}

func (m *MetricsUC) Find(ctx context.Context, name string, mType models.MetricType) (*models.Metric, error) {
	return m.Storage.Get(ctx, name, mType)
}

func (m *MetricsUC) Update(ctx context.Context, metric models.Metric) error {
	return m.Storage.Update(ctx, metric)
}

func (m *MetricsUC) Updates(ctx context.Context, metrics []models.Metric) error {
	return m.Storage.Updates(ctx, metrics)
}

func (m *MetricsUC) GetAll(ctx context.Context) ([]models.Metric, error) {
	return m.Storage.GetAll(ctx)
}

func NewMetricsUC(storage repository.Repository, config *config.ServerConfig) *MetricsUC {
	return &MetricsUC{Storage: storage, config: config}
}
