package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
)

func TestMemStorage_Update(t *testing.T) {
	memStorage := NewMemStorage()

	repository.TestUpdate(context.Background(), t, memStorage)
}

func TestMemStorage_Get(t *testing.T) {
	storage := NewMemStorage()
	require.NoError(t, storage.Update(context.Background(), "Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update(context.Background(), "PollCount", models.Counter, "1"))

	repository.TestGet(context.Background(), t, storage)
}

func TestMemStorage_GetAll(t *testing.T) {
	storage := NewMemStorage()

	require.NoError(t, storage.Update(context.Background(), "Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update(context.Background(), "PollCount", models.Counter, "1"))

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
