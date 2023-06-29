package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
		want Metrics
	}{
		{
			name: "Create metrics",
			want: Metrics{
				"Alloc":         NewMetric(Gauge),
				"BuckHashSys":   NewMetric(Gauge),
				"Frees":         NewMetric(Gauge),
				"GCCPUFraction": NewMetric(Gauge),
				"GCSys":         NewMetric(Gauge),
				"HeapAlloc":     NewMetric(Gauge),
				"HeapIdle":      NewMetric(Gauge),
				"HeapInuse":     NewMetric(Gauge),
				"HeapObjects":   NewMetric(Gauge),
				"HeapReleased":  NewMetric(Gauge),
				"HeapSys":       NewMetric(Gauge),
				"LastGC":        NewMetric(Gauge),
				"Lookups":       NewMetric(Gauge),
				"MCacheInuse":   NewMetric(Gauge),
				"MCacheSys":     NewMetric(Gauge),
				"MSpanInuse":    NewMetric(Gauge),
				"MSpanSys":      NewMetric(Gauge),
				"Mallocs":       NewMetric(Gauge),
				"NextGC":        NewMetric(Gauge),
				"NumForcedGC":   NewMetric(Gauge),
				"NumGC":         NewMetric(Gauge),
				"OtherSys":      NewMetric(Gauge),
				"PauseTotalNs":  NewMetric(Gauge),
				"StackInuse":    NewMetric(Gauge),
				"StackSys":      NewMetric(Gauge),
				"Sys":           NewMetric(Gauge),
				"TotalAlloc":    NewMetric(Gauge),
				"RandomValue":   NewMetric(Gauge),
				"PollCount":     NewMetric(Counter),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMetrics()
			for key, metric := range got {
				got[key] = NewMetric(metric.ValueType())
			}

			assert.ObjectsAreEqualValues(tt.want, got)
		})
	}
}
