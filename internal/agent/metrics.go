package agent

import (
	"go-metricscol/internal/models"
	"math/rand"
	"runtime"
)

func UpdateMetrics(metrics models.Metrics) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	metrics.UpdateGauge("Alloc", float64(stats.Alloc))
	metrics.UpdateGauge("BuckHashSys", float64(stats.BuckHashSys))
	metrics.UpdateGauge("BuckHashSys", float64(stats.BuckHashSys))
	metrics.UpdateGauge("Frees", float64(stats.Frees))
	metrics.UpdateGauge("GCCPUFraction", stats.GCCPUFraction)
	metrics.UpdateGauge("GCSys", float64(stats.GCSys))
	metrics.UpdateGauge("HeapAlloc", float64(stats.HeapAlloc))
	metrics.UpdateGauge("HeapIdle", float64(stats.HeapIdle))
	metrics.UpdateGauge("HeapInuse", float64(stats.HeapInuse))
	metrics.UpdateGauge("HeapObjects", float64(stats.HeapObjects))
	metrics.UpdateGauge("HeapReleased", float64(stats.HeapReleased))
	metrics.UpdateGauge("HeapSys", float64(stats.HeapSys))
	metrics.UpdateGauge("LastGC", float64(stats.LastGC))
	metrics.UpdateGauge("Lookups", float64(stats.Lookups))
	metrics.UpdateGauge("MCacheInuse", float64(stats.MCacheInuse))
	metrics.UpdateGauge("MCacheSys", float64(stats.MCacheSys))
	metrics.UpdateGauge("MSpanInuse", float64(stats.MSpanInuse))
	metrics.UpdateGauge("MSpanSys", float64(stats.MSpanSys))
	metrics.UpdateGauge("Mallocs", float64(stats.Mallocs))
	metrics.UpdateGauge("NextGC", float64(stats.NextGC))
	metrics.UpdateGauge("NumForcedGC", float64(stats.NumForcedGC))
	metrics.UpdateGauge("NumGC", float64(stats.NumGC))
	metrics.UpdateGauge("OtherSys", float64(stats.OtherSys))
	metrics.UpdateGauge("PauseTotalNs", float64(stats.PauseTotalNs))
	metrics.UpdateGauge("StackInuse", float64(stats.StackInuse))
	metrics.UpdateGauge("StackSys", float64(stats.StackSys))
	metrics.UpdateGauge("Sys", float64(stats.Sys))
	metrics.UpdateGauge("TotalAlloc", float64(stats.TotalAlloc))
	metrics.UpdateGauge("RandomValue", rand.Float64())
	metrics.UpdateCounter("PollCount")
}

func CreateMetrics() models.Metrics {
	metrics := models.NewMetrics()
	UpdateMetrics(metrics)

	return metrics
}
