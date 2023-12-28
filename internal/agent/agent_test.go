package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func TestUpdateMetrics(t *testing.T) {
	metrics := memory.NewMetrics()

	metricsMustBeUpdated := []string{"BuckHashSys", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapSys", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "OtherSys", "StackInuse", "StackSys", "Sys", "TotalAlloc", "RandomValue", "PollCount", "Alloc"}
	t.Run("UpdateMetrics", func(t *testing.T) {
		assert.NoError(t, UpdateMetrics(&metrics))

		for key, metric := range metrics.Collection {
			if contains(metricsMustBeUpdated, key) {
				assert.NotEqual(t, metric.StringValue(), "0")
			}
		}
	})
}

func BenchmarkUpdateMetrics(b *testing.B) {
	metrics := memory.NewMetrics()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assert.NoError(b, UpdateMetrics(&metrics))
	}
}

func BenchmarkCollectAdditionalMetrics(b *testing.B) {
	metrics := memory.NewMetrics()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assert.NoError(b, CollectAdditionalMetrics(&metrics))
	}
}

func TestUpdatePollCount(t *testing.T) {
	metrics := memory.NewMetrics()
	assert.NoError(t, UpdateMetrics(&metrics))
	assert.NoError(t, UpdateMetrics(&metrics))
	assert.NoError(t, UpdateMetrics(&metrics))
	assert.NoError(t, UpdateMetrics(&metrics))
	assert.NoError(t, UpdateMetrics(&metrics))

	pollCount, err := metrics.Get("PollCount", models.Counter)
	assert.EqualValues(t, nil, err)
	assert.Equal(t, pollCount.StringValue(), "5")
}
