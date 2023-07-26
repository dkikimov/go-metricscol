package repository

import (
	"go-metricscol/internal/models"
)

type Repository interface {
	Update(name string, valueType models.MetricType, value string) error
	Get(key string, valueType models.MetricType) (*models.Metric, error)
	UpdateWithStruct(metric *models.Metric) error
	GetAll() ([]models.Metric, error)

	//TODO: Убрал маршалинг, тк не нужен постгресу. Нужно подумать как требовать маршалинг от MemStorage
	//json.Marshaler
	//json.Unmarshaler
}
