package memory

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"strconv"
)

type MemStorage struct {
	metrics models.Metrics
}

func (memStorage *MemStorage) GetAll() map[string]models.Metric {
	return memStorage.metrics
}

func (memStorage *MemStorage) Get(name string, valueType models.MetricType) (models.Metric, apierror.APIError) {
	return memStorage.metrics.Get(name, valueType)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: models.Metrics{}}
}

func (memStorage *MemStorage) Update(name string, valueType models.MetricType, value string) apierror.APIError {
	switch valueType {
	case models.GaugeType:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return apierror.NumberParse
		}
		memStorage.metrics.Update(name, models.GaugeType, floatVal)
	case models.CounterType:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return apierror.NumberParse
		}
		memStorage.metrics.Update(name, models.CounterType, intVal)
	default:
		return apierror.UnknownMetricType
	}

	return apierror.NoError
}
