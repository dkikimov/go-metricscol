package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-metricscol/internal/config"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/utils"
)

var testMetric = models.Metric{
	Name:  "test",
	MType: models.Gauge,
	Value: utils.Ptr(float64(1)),
}

func TestServer_enableSavingToDisk(t *testing.T) {
	file, err := os.CreateTemp("", "test_restoreFromDisk")

	defer func() {
		require.NoError(t, os.Remove(file.Name()))
	}()

	require.NoError(t, file.Close())
	require.NoError(t, err)

	storeInterval := 2 * time.Second

	cfg, err := config.NewServerConfig("127.0.0.1:8080", models.Duration{Duration: storeInterval}, file.Name(), false, "", "", "", "")
	require.NoError(t, err)

	server := NewServer(cfg, memory.NewMemStorage(), nil)
	require.NoError(t, server.Repo.UpdateWithStruct(context.Background(), &testMetric))
	require.NoError(t, err)

	t.Run("Enable saving to disk", func(t *testing.T) {
		go server.enableSavingToDisk(context.Background())

		time.Sleep(storeInterval - 1*time.Second)

		bytes, err := os.ReadFile(file.Name())
		require.NoError(t, err)
		assert.Equal(t, []byte{}, bytes)

		time.Sleep(1500 * time.Millisecond)

		savedStorage := memory.NewMemStorage()
		require.NoError(t, savedStorage.RestoreFromDisk(file.Name()))

		assert.Equal(t, server.Repo, savedStorage)
	})
}
