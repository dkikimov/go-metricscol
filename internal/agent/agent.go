package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
)

func SendMetricsToServer(addr string, m *memory.Metrics, hashKey string) error {
	for _, metric := range m.Collection {
		postURL := url.URL{
			Scheme: "http",
			Host:   addr,
			Path:   "/update/",
		}

		log.Println(postURL.String())

		metric.Hash = metric.HashValue(hashKey)
		marshal, err := json.Marshal(metric)
		if err != nil {
			return errors.New("couldn't marshal metric")
		}

		resp, err := http.Post(postURL.String(), "application/json", bytes.NewReader(marshal))

		if err != nil {
			return fmt.Errorf("couldn't post url %s", postURL.String())
		}

		if err := resp.Body.Close(); err != nil {
			return errors.New("couldn't close response body")
		}
	}
	m.ResetPollCount()

	return nil
}

func UpdateMetrics(metrics *memory.Metrics) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	// TODO: Стоит ли хенлдить ошибки?
	metrics.Update("Alloc", models.Gauge, float64(stats.Alloc))
	metrics.Update("BuckHashSys", models.Gauge, float64(stats.BuckHashSys))
	metrics.Update("BuckHashSys", models.Gauge, float64(stats.BuckHashSys))
	metrics.Update("Frees", models.Gauge, float64(stats.Frees))
	metrics.Update("GCCPUFraction", models.Gauge, stats.GCCPUFraction)
	metrics.Update("GCSys", models.Gauge, float64(stats.GCSys))
	metrics.Update("HeapAlloc", models.Gauge, float64(stats.HeapAlloc))
	metrics.Update("HeapIdle", models.Gauge, float64(stats.HeapIdle))
	metrics.Update("HeapInuse", models.Gauge, float64(stats.HeapInuse))
	metrics.Update("HeapObjects", models.Gauge, float64(stats.HeapObjects))
	metrics.Update("HeapReleased", models.Gauge, float64(stats.HeapReleased))
	metrics.Update("HeapSys", models.Gauge, float64(stats.HeapSys))
	metrics.Update("LastGC", models.Gauge, float64(stats.LastGC))
	metrics.Update("Lookups", models.Gauge, float64(stats.Lookups))
	metrics.Update("MCacheInuse", models.Gauge, float64(stats.MCacheInuse))
	metrics.Update("MCacheSys", models.Gauge, float64(stats.MCacheSys))
	metrics.Update("MSpanInuse", models.Gauge, float64(stats.MSpanInuse))
	metrics.Update("MSpanSys", models.Gauge, float64(stats.MSpanSys))
	metrics.Update("Mallocs", models.Gauge, float64(stats.Mallocs))
	metrics.Update("NextGC", models.Gauge, float64(stats.NextGC))
	metrics.Update("NumForcedGC", models.Gauge, float64(stats.NumForcedGC))
	metrics.Update("NumGC", models.Gauge, float64(stats.NumGC))
	metrics.Update("OtherSys", models.Gauge, float64(stats.OtherSys))
	metrics.Update("PauseTotalNs", models.Gauge, float64(stats.PauseTotalNs))
	metrics.Update("StackInuse", models.Gauge, float64(stats.StackInuse))
	metrics.Update("StackSys", models.Gauge, float64(stats.StackSys))
	metrics.Update("Sys", models.Gauge, float64(stats.Sys))
	metrics.Update("TotalAlloc", models.Gauge, float64(stats.TotalAlloc))
	metrics.Update("RandomValue", models.Gauge, rand.Float64())
	metrics.Update("PollCount", models.Counter, 1)
}
