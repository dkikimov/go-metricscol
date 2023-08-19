package memory

import (
	"context"
	"encoding/json"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"sort"
	"strconv"
)

type MemStorage struct {
	metrics Metrics
}

func (memStorage *MemStorage) SupportsSavingToDisk() bool {
	return true
}

func (memStorage *MemStorage) SupportsTx() bool {
	return false
}

func (memStorage *MemStorage) Updates(_ context.Context, metrics []models.Metric) error {
	for _, metric := range metrics {
		err := memStorage.metrics.UpdateWithStruct(&metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (memStorage *MemStorage) UnmarshalJSON(bytes []byte) error {
	// TODO: Подумать
	memStorage.metrics.mu.Lock()
	defer memStorage.metrics.mu.Unlock()

	return json.Unmarshal(bytes, &memStorage.metrics.Collection)
}

func (memStorage *MemStorage) MarshalJSON() ([]byte, error) {
	memStorage.metrics.mu.RLock()
	defer memStorage.metrics.mu.RUnlock()

	return json.Marshal(memStorage.metrics.Collection)
}

func (memStorage *MemStorage) UpdateWithStruct(_ context.Context, metric *models.Metric) error {
	return memStorage.metrics.UpdateWithStruct(metric)
}

func (memStorage *MemStorage) GetAll(context.Context) ([]models.Metric, error) {
	all := memStorage.metrics.GetAll()

	sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })

	return all, nil
}

func (memStorage *MemStorage) Get(_ context.Context, key string, valueType models.MetricType) (*models.Metric, error) {
	result, err := memStorage.metrics.Get(key, valueType)
	return result, err
}

func (memStorage *MemStorage) Update(_ context.Context, name string, valueType models.MetricType, value string) error {
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

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: NewMetrics()}
}
