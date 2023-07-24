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
	metrics models.Metrics
	mu      sync.Mutex
	config  *Config
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

	all := memStorage.metrics.GetAll()
	for idx, value := range all {
		all[idx].Hash = value.HashValue(memStorage.config.HashKey)
	}

	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })

	return all
}

func (memStorage *MemStorage) Get(key string, valueType models.MetricType) (*models.Metric, error) {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	result, err := memStorage.metrics.Get(key, valueType)
	if err == nil {
		result.Hash = result.HashValue(memStorage.config.HashKey)
	}

	return result, err
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

func NewMemStorage(hashKey string) *MemStorage {
	return &MemStorage{metrics: models.NewMetrics(), config: NewConfig(hashKey)}
}
