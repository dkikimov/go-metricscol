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
	metrics := models.NewMetrics()

	metricsMustBeUpdated := []string{"BuckHashSys", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapSys", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "OtherSys", "StackInuse", "StackSys", "Sys", "TotalAlloc", "RandomValue", "PollCount", "Alloc"}
	t.Run("UpdateMetrics", func(t *testing.T) {
		UpdateMetrics(metrics)

		for key, metric := range metrics {
			if contains(metricsMustBeUpdated, key) {
				assert.NotEqual(t, metric.StringValue(), "0")
			}
		}
	})
}

func TestUpdatePollCount(t *testing.T) {
	metrics := models.NewMetrics()
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)
	UpdateMetrics(metrics)

	assert.Equal(t, metrics["PollCount"].StringValue(), "5")
}
