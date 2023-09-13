package agent

import (
	"testing"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"

	"github.com/stretchr/testify/assert"
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
		UpdateMetrics(&metrics)

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
		UpdateMetrics(&metrics)
	}
}

func BenchmarkCollectAdditionalMetrics(b *testing.B) {
	metrics := memory.NewMetrics()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CollectAdditionalMetrics(&metrics)
	}
}

func TestUpdatePollCount(t *testing.T) {
	metrics := memory.NewMetrics()
	UpdateMetrics(&metrics)
	UpdateMetrics(&metrics)
	UpdateMetrics(&metrics)
	UpdateMetrics(&metrics)
	UpdateMetrics(&metrics)

	pollCount, err := metrics.Get("PollCount", models.Counter)
	assert.EqualValues(t, nil, err)
	assert.Equal(t, pollCount.StringValue(), "5")
}
