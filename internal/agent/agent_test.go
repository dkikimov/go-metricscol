package agent

import (
	"github.com/stretchr/testify/assert"
	"go-metricscol/internal/models"
	"testing"
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
	metrics := models.Metrics{}

	metricsMustBeUpdated := []string{"BuckHashSys", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapSys", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "OtherSys", "StackInuse", "StackSys", "Sys", "TotalAlloc", "RandomValue", "PollCount", "Alloc"}
	t.Run("UpdateMetrics", func(t *testing.T) {
		UpdateMetrics(metrics)

		for key, metric := range metrics {
			if contains(metricsMustBeUpdated, key) {
				assert.NotEqual(t, metric.GetStringValue(), "0")
			}
		}
	})
}

func TestUpdatePollCount(t *testing.T) {
	metrics := models.Metrics{}
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)

	pollCount, err := metrics.Get("PollCount", models.CounterType)
	assert.EqualValues(t, nil, err)
	assert.Equal(t, pollCount.GetStringValue(), "5")
}
