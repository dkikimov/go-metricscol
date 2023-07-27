package server

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/utils"
	"os"
	"testing"
	"time"
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
	config := NewConfig("127.0.0.1:8080", storeInterval, file.Name(), false, "", "")

	server, err := NewServer(config)
	require.NoError(t, server.Repository.UpdateWithStruct(&testMetric))
	require.NoError(t, err)

	t.Run("Enable saving to disk", func(t *testing.T) {
		go server.enableSavingToDisk()

		time.Sleep(storeInterval - 1*time.Second)

		bytes, err := os.ReadFile(file.Name())
		require.NoError(t, err)
		assert.Equal(t, []byte{}, bytes)

		time.Sleep(1500 * time.Millisecond)

		bytes, err = os.ReadFile(file.Name())
		require.NoError(t, err)

		savedStorage := memory.NewMemStorage()
		err = json.Unmarshal(bytes, savedStorage)
		require.NoError(t, err)

		assert.Equal(t, server.Repository, savedStorage)
	})
}

func TestServer_restoreFromDisk(t *testing.T) {
	file, err := os.CreateTemp("", "test_restoreFromDisk")

	require.NoError(t, os.Remove(file.Name()))

	require.NoError(t, file.Close())
	require.NoError(t, err)

	config := NewConfig("127.0.0.1:8080", 5*time.Second, file.Name(), false, "", "")
	storage := memory.NewMemStorage()

	require.NoError(t, storage.UpdateWithStruct(&testMetric))

	server, err := NewServer(config)
	require.NoError(t, err)

	require.NoError(t, server.saveToDisk())

	t.Run("Restore from disk", func(t *testing.T) {
		newServer, err := NewServer(config)
		require.NoError(t, err)
		require.NoError(t, newServer.restoreFromDisk())

		assert.Equal(t, server, newServer)
	})
}

func TestServer_saveToDisk(t *testing.T) {
	type fields struct {
		Config     *Config
		Repository repository.Repository
	}

	file, err := os.CreateTemp("", "test_saveToDisk")

	require.NoError(t, os.Remove(file.Name()))

	require.NoError(t, file.Close())
	require.NoError(t, err)

	config := NewConfig("127.0.0.1:8080", 5*time.Second, file.Name(), false, "", "")
	storage := memory.NewMemStorage()

	require.NoError(t, storage.UpdateWithStruct(&testMetric))

	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "save to disk",
			fields: fields{
				Config:     config,
				Repository: storage,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{
				Config:     tt.fields.Config,
				Repository: tt.fields.Repository,
			}

			err := s.saveToDisk()

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			bytes, err := os.ReadFile(file.Name())
			require.NoError(t, err)

			savedStorage := memory.NewMemStorage()
			err = json.Unmarshal(bytes, savedStorage)
			require.NoError(t, err)

			assert.Equal(t, storage, savedStorage)

		})
	}
}
