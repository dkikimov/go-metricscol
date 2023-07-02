package agent

import (
	"errors"
	"fmt"
	"go-metricscol/internal/models"
	"log"
	"math/rand"
	"net/http"
	"runtime"
)

func SendMetricsToServer(addr string, m models.Metrics) error {
	for _, metric := range m {
		postURL := fmt.Sprintf("%s/update/%s/%s/%s", addr, metric.GetType(), metric.GetName(), metric.GetStringValue())
		log.Println(postURL)
		resp, err := http.Post(postURL, "text/plain", nil)

		if err != nil {
			return fmt.Errorf("couldn't post url %s", postURL)
		}

		if err := resp.Body.Close(); err != nil {
			return errors.New("couldn't close response body")
		}
	}
	m.ResetPollCount()

	return nil
}

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
