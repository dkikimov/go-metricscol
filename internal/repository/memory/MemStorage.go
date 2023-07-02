package memory

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"sort"
	"strconv"
)

type MemStorage struct {
	metrics models.Metrics
}

func (memStorage *MemStorage) GetAll() []models.Metric {
	kv := make([]models.Metric, 0, len(memStorage.metrics))
	for _, value := range memStorage.metrics {
		kv = append(kv, value)
	}

	sort.Slice(kv, func(i, j int) bool { return kv[i].GetName() < kv[j].GetName() })

	return kv
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
