package agent

import (
	"go-metricscol/internal/models"
	"math/rand"
	"runtime"
)

func UpdateMetrics(metrics models.Metrics) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	metrics.Update("Alloc", models.GaugeType, float64(stats.Alloc))
	metrics.Update("BuckHashSys", models.GaugeType, float64(stats.BuckHashSys))
	metrics.Update("BuckHashSys", models.GaugeType, float64(stats.BuckHashSys))
	metrics.Update("Frees", models.GaugeType, float64(stats.Frees))
	metrics.Update("GCCPUFraction", models.GaugeType, stats.GCCPUFraction)
	metrics.Update("GCSys", models.GaugeType, float64(stats.GCSys))
	metrics.Update("HeapAlloc", models.GaugeType, float64(stats.HeapAlloc))
	metrics.Update("HeapIdle", models.GaugeType, float64(stats.HeapIdle))
	metrics.Update("HeapInuse", models.GaugeType, float64(stats.HeapInuse))
	metrics.Update("HeapObjects", models.GaugeType, float64(stats.HeapObjects))
	metrics.Update("HeapReleased", models.GaugeType, float64(stats.HeapReleased))
	metrics.Update("HeapSys", models.GaugeType, float64(stats.HeapSys))
	metrics.Update("LastGC", models.GaugeType, float64(stats.LastGC))
	metrics.Update("Lookups", models.GaugeType, float64(stats.Lookups))
	metrics.Update("MCacheInuse", models.GaugeType, float64(stats.MCacheInuse))
	metrics.Update("MCacheSys", models.GaugeType, float64(stats.MCacheSys))
	metrics.Update("MSpanInuse", models.GaugeType, float64(stats.MSpanInuse))
	metrics.Update("MSpanSys", models.GaugeType, float64(stats.MSpanSys))
	metrics.Update("Mallocs", models.GaugeType, float64(stats.Mallocs))
	metrics.Update("NextGC", models.GaugeType, float64(stats.NextGC))
	metrics.Update("NumForcedGC", models.GaugeType, float64(stats.NumForcedGC))
	metrics.Update("NumGC", models.GaugeType, float64(stats.NumGC))
	metrics.Update("OtherSys", models.GaugeType, float64(stats.OtherSys))
	metrics.Update("PauseTotalNs", models.GaugeType, float64(stats.PauseTotalNs))
	metrics.Update("StackInuse", models.GaugeType, float64(stats.StackInuse))
	metrics.Update("StackSys", models.GaugeType, float64(stats.StackSys))
	metrics.Update("Sys", models.GaugeType, float64(stats.Sys))
	metrics.Update("TotalAlloc", models.GaugeType, float64(stats.TotalAlloc))
	metrics.Update("RandomValue", models.GaugeType, rand.Float64())
	metrics.Update("PollCount", models.CounterType, 1)
}
