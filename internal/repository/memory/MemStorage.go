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

func (memStorage *MemStorage) Get(key string, valueType models.MetricType) (models.Metric, apierror.APIError) {
	metric, ok := memStorage.metrics[key]

	// TODO: Возможно стоит Добавить поддержку метрик с одинаковым названием и разными типами
	if !ok || metric.ValueType() != valueType {
		return models.Metric{}, apierror.NotFound
	}

	return metric, apierror.NoError
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: models.Metrics{}}
}

func (memStorage *MemStorage) Update(key string, value string, valueType models.MetricType) apierror.APIError {
	_, ok := memStorage.metrics[key]
	if !ok {
		memStorage.metrics[key] = models.NewMetric(valueType)
	}
	switch valueType {
	case models.Gauge:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return apierror.NumberParse
		}
		memStorage.metrics.UpdateGauge(key, floatVal)
	case models.Counter:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return apierror.NumberParse
		}
		memStorage.metrics.UpdateCounter(key, intVal)
	}

	return apierror.NoError
}
