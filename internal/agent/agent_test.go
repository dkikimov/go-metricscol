package agent

import (
	"github.com/stretchr/testify/assert"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
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
	metrics := models.MetricsMap{}

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
	metrics := models.MetricsMap{}
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)

	pollCount, err := metrics.Get("PollCount", models.Counter)
	assert.EqualValues(t, apierror.NoError, err)
	assert.Equal(t, pollCount.GetStringValue(), "5")
}
