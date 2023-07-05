package memory

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"sort"
	"strconv"
)

type MemStorage struct {
	metrics models.MetricsMap
}

func (memStorage *MemStorage) UpdateWithStruct(metric *models.Metric) apierror.APIError {
	return memStorage.metrics.UpdateWithStruct(metric)
}

func (memStorage *MemStorage) GetAll() []models.Metric {
	kv := make([]models.Metric, 0, len(memStorage.metrics))
	for _, value := range memStorage.metrics {
		kv = append(kv, value)
	}

	sort.Slice(kv, func(i, j int) bool { return kv[i].Name < kv[j].Name })

	return kv
}

func (memStorage *MemStorage) Get(name string, valueType models.MetricType) (*models.Metric, apierror.APIError) {
	return memStorage.metrics.Get(name, valueType)
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: models.MetricsMap{}}
}

func (memStorage *MemStorage) Update(name string, valueType models.MetricType, value string) apierror.APIError {
	switch valueType {
	case models.Gauge:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return apierror.NumberParse
		}
		memStorage.metrics.Update(name, models.Gauge, floatVal)
	case models.Counter:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return apierror.NumberParse
		}
		memStorage.metrics.Update(name, models.Counter, intVal)
	default:
		return apierror.UnknownMetricType
	}

	return apierror.NoError
}
