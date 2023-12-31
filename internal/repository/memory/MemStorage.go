package memory

import (
	"context"
	"encoding/json"
	"sort"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
)

// MemStorage is a metrics in-memory storage which implements Repository interface.
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

func (memStorage *MemStorage) Update(ctx context.Context, metric models.Metric) error {
	switch metric.MType {
	case models.Gauge:
		// TODO: update signature
		return memStorage.metrics.Update(metric.Name, models.Gauge, *metric.Value)
	case models.Counter:
		return memStorage.metrics.Update(metric.Name, models.Gauge, *metric.Delta)
	default:
		return apierror.UnknownMetricType
	}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: NewMetrics()}
}
