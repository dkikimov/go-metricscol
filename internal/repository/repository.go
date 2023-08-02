package repository

import (
	"encoding/json"
	"go-metricscol/internal/models"
)

type Repository interface {
	Update(name string, valueType models.MetricType, value string) error
	Updates(metrics []models.Metric) error
	UpdateWithStruct(metric *models.Metric) error
	Get(key string, valueType models.MetricType) (*models.Metric, error)
	GetAll() ([]models.Metric, error)
	SupportsTx() bool

	json.Marshaler
	json.Unmarshaler
}
