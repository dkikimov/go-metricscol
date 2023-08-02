package memory

import (
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"testing"
)

func TestMemStorage_Update(t *testing.T) {
	memStorage := NewMemStorage()

	repository.TestUpdate(t, memStorage)
}

func TestMemStorage_Get(t *testing.T) {
	storage := NewMemStorage()
	require.NoError(t, storage.Update("Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update("PollCount", models.Counter, "1"))

	repository.TestGet(t, storage)
}

func TestMemStorage_GetAll(t *testing.T) {
	storage := NewMemStorage()

	require.NoError(t, storage.Update("Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update("PollCount", models.Counter, "1"))

	repository.TestGetAll(t, storage)
}

func TestMemStorage_UpdateWithStruct(t *testing.T) {
	storage := NewMemStorage()

	repository.TestUpdateWithStruct(t, storage)
}

func TestMemStorage_Updates(t *testing.T) {
	storage := NewMemStorage()

	repository.TestUpdates(t, storage)
}
