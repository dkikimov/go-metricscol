package memory

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"strconv"
)

type MemStorage struct {
	metrics models.Metrics
}

func (memStorage *MemStorage) GetString(key string, valueType models.MetricType) (string, apierror.APIError) {
	metric, ok := memStorage.metrics[key]

	// TODO: Возможно стоит Добавить поддержку метрик с одинаковым названием и разными типами
	if !ok || metric.ValueType() != valueType {
		return "", apierror.NotFound
	}

	return metric.StringValue(), apierror.NoError
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
