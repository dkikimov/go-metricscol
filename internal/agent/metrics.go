package agent

import (
	"go-metricscol/internal/models"
	"math/rand"
	"runtime"
)

func UpdateMetrics(metrics *models.Metrics) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	metrics.Alloc = models.Gauge(stats.Alloc)
	metrics.BuckHashSys = models.Gauge(stats.BuckHashSys)
	metrics.Frees = models.Gauge(stats.Frees)
	metrics.GCCPUFraction = models.Gauge(stats.GCCPUFraction)
	metrics.GCSys = models.Gauge(stats.GCSys)
	metrics.HeapAlloc = models.Gauge(stats.HeapAlloc)
	metrics.HeapIdle = models.Gauge(stats.HeapIdle)
	metrics.HeapInuse = models.Gauge(stats.HeapInuse)
	metrics.HeapObjects = models.Gauge(stats.HeapObjects)
	metrics.HeapReleased = models.Gauge(stats.HeapReleased)
	metrics.HeapSys = models.Gauge(stats.HeapSys)
	metrics.LastGC = models.Gauge(stats.LastGC)
	metrics.Lookups = models.Gauge(stats.Lookups)
	metrics.MCacheInuse = models.Gauge(stats.MCacheInuse)
	metrics.MCacheSys = models.Gauge(stats.MCacheSys)
	metrics.MSpanInuse = models.Gauge(stats.MSpanInuse)
	metrics.MSpanSys = models.Gauge(stats.MSpanSys)
	metrics.Mallocs = models.Gauge(stats.Mallocs)
	metrics.NextGC = models.Gauge(stats.NextGC)
	metrics.NumForcedGC = models.Gauge(stats.NumForcedGC)
	metrics.NumGC = models.Gauge(stats.NumGC)
	metrics.OtherSys = models.Gauge(stats.OtherSys)
	metrics.PauseTotalNs = models.Gauge(stats.PauseTotalNs)
	metrics.StackInuse = models.Gauge(stats.StackInuse)
	metrics.StackSys = models.Gauge(stats.StackSys)
	metrics.Sys = models.Gauge(stats.Sys)
	metrics.TotalAlloc = models.Gauge(stats.TotalAlloc)

	metrics.RandomValue = models.Gauge(rand.Float64())
	metrics.PollCount++
}

func CreateMetrics() *models.Metrics {
	metrics := new(models.Metrics)
	UpdateMetrics(metrics)

	return metrics
}
