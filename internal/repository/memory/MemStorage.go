package memory

import (
	"encoding/json"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"sort"
	"strconv"
	"sync"
)

type MemStorage struct {
	metrics models.MetricsMap
	mu      sync.RWMutex
}

func (memStorage *MemStorage) UnmarshalJSON(bytes []byte) error {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	return json.Unmarshal(bytes, &memStorage.metrics)
}

func (memStorage *MemStorage) MarshalJSON() ([]byte, error) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	return json.Marshal(memStorage.metrics)
}

func (memStorage *MemStorage) UpdateWithStruct(metric *models.Metric) error {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	return memStorage.metrics.UpdateWithStruct(metric)
}

func (memStorage *MemStorage) GetAll() []models.Metric {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	kv := make([]models.Metric, 0, len(memStorage.metrics))
	for _, value := range memStorage.metrics {
		kv = append(kv, value)
	}

	sort.Slice(kv, func(i, j int) bool { return kv[i].Name < kv[j].Name })

	return kv
}

func (memStorage *MemStorage) Get(key string, valueType models.MetricType) (*models.Metric, error) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	result, err := memStorage.metrics.Get(key, valueType)
	return result, err
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: models.MetricsMap{}}
}

func (memStorage *MemStorage) Update(name string, valueType models.MetricType, value string) error {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	switch valueType {
	case models.Gauge:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return apierror.NumberParse
		}
		return memStorage.metrics.Update(name, models.Gauge, floatVal)
	case models.Counter:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return apierror.NumberParse
		}
		return memStorage.metrics.Update(name, models.Counter, intVal)
	default:
		return apierror.UnknownMetricType
	}
}
