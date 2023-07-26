package memory

import (
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"os"
	"testing"
)

func TestMemStorage_Update(t *testing.T) {
	memStorage := NewMemStorage("")

	repository.TestUpdate(t, memStorage)
}

func TestMemStorage_Get(t *testing.T) {
	storage := NewMemStorage("")
	require.NoError(t, storage.Update("Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update("PollCount", models.Counter, "2"))

	repository.TestGet(t, storage)
}

func TestMemStorage_GetAll(t *testing.T) {
	storage := NewMemStorage("")

	require.NoError(t, storage.Update("Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update("PollCount", models.Counter, "2"))
	require.NoError(t, os.Setenv("KEY", ""))

	repository.TestGetAll(t, storage)
}
