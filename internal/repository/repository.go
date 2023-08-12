package repository

import (
	"context"
	"go-metricscol/internal/models"
)

type Repository interface {
	Update(ctx context.Context, name string, valueType models.MetricType, value string) error
	Updates(ctx context.Context, metrics []models.Metric) error
	UpdateWithStruct(ctx context.Context, metric *models.Metric) error
	Get(ctx context.Context, key string, valueType models.MetricType) (*models.Metric, error)
	GetAll(ctx context.Context) ([]models.Metric, error)
	SupportsTx() bool
	SupportsSavingToDisk() bool

	//json.Marshaler
	//json.Unmarshaler
}
