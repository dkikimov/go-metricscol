package memory

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/utils"
)

func TestMemStorage_Update(t *testing.T) {
	memStorage := NewMemStorage()

	repository.TestUpdate(context.Background(), t, memStorage)
}

func TestMemStorage_Get(t *testing.T) {
	storage := NewMemStorage()

	require.NoError(t, storage.Update(context.Background(), models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(101.42),
	}))

	require.NoError(t, storage.Update(context.Background(), models.Metric{
		Name:  "PollCount",
		MType: models.Counter,
		Delta: utils.Ptr(int64(1)),
	}))

	repository.TestGet(context.Background(), t, storage)
}

func TestMemStorage_GetAll(t *testing.T) {
	storage := NewMemStorage()

	require.NoError(t, storage.Update(context.Background(), models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(101.42),
	}))

	require.NoError(t, storage.Update(context.Background(), models.Metric{
		Name:  "PollCount",
		MType: models.Counter,
		Delta: utils.Ptr(int64(1)),
	}))

	repository.TestGetAll(context.Background(), t, storage)
}

func TestMemStorage_UpdateWithStruct(t *testing.T) {
	storage := NewMemStorage()

	repository.TestUpdateWithStruct(context.Background(), t, storage)
}

func TestMemStorage_Updates(t *testing.T) {
	storage := NewMemStorage()

	repository.TestUpdates(context.Background(), t, storage)
}

var testMetric = models.Metric{
	Name:  "test",
	MType: models.Gauge,
	Value: utils.Ptr(float64(1)),
}

func TestMemStorage_SaveAndRestoreFromDisk(t *testing.T) {
	file, err := os.CreateTemp("", "test_RestoreFromDisk")

	require.NoError(t, os.Remove(file.Name()))

	require.NoError(t, file.Close())
	require.NoError(t, err)

	storage := NewMemStorage()

	require.NoError(t, storage.UpdateWithStruct(context.Background(), &testMetric))
	require.NoError(t, storage.SaveToDisk(file.Name()))

	t.Run("Restore from disk", func(t *testing.T) {
		newStorage := NewMemStorage()

		require.NoError(t, newStorage.RestoreFromDisk(file.Name()))
		assert.Equal(t, storage, newStorage)
	})
}
