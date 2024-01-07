package memory

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sort"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
)

// MemStorage is a metrics in-memory storage which implements Repository interface.
type MemStorage struct {
	metrics Metrics
}

func (memStorage *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (memStorage *MemStorage) RestoreFromDisk(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&memStorage.metrics.Collection); err != nil {
		return err
	}

	return nil
}

func (memStorage *MemStorage) SaveToDisk(filePath string) error {
	log.Printf("saving to disk")

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(memStorage.metrics.Collection); err != nil {
		return err
	}

	return nil
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
		return memStorage.metrics.Update(metric.Name, models.Counter, *metric.Delta)
	default:
		return apierror.UnknownMetricType
	}
}

func NewMemStorage() *MemStorage {
	return &MemStorage{metrics: NewMetrics()}
}
