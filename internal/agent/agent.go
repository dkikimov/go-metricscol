package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"golang.org/x/sync/errgroup"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
)

func SendMetricsToServer(cfg *Config, m *memory.Metrics) error {
	jobCh := make(chan bool)
	g := errgroup.Group{}

	for i := 0; i < cfg.RateLimit; i++ {
		g.Go(func() error {
			for range jobCh {
				return makeRequest(cfg, m)
			}
			return nil
		})
	}

	jobCh <- true
	close(jobCh)

	if err := g.Wait(); err != nil {
		return err
	}

	m.ResetPollCount()
	return nil
}

func makeRequest(cfg *Config, m *memory.Metrics) error {
	postURL := url.URL{
		Scheme: "http",
		Host:   cfg.Address,
		Path:   "/updates/",
	}

	metrics := make([]models.Metric, 0, len(m.Collection))
	for _, value := range m.Collection {
		value.Hash = value.HashValue(cfg.HashKey)
		metrics = append(metrics, value)
	}

	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		return errors.New("couldn't marshal metrics")
	}

	gzipMetrics := bytes.NewBuffer([]byte{})
	w := gzip.NewWriter(gzipMetrics)
	_, err = w.Write(jsonMetrics)
	if err != nil {
		return fmt.Errorf("couldn't gzip metrics with error: %s", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("couldn't close gzip writer with error: %s", err)
	}

	request, err := http.NewRequest(http.MethodPost, postURL.String(), gzipMetrics)
	if err != nil {
		return fmt.Errorf("couldn't create request with error: %s", err)
	}
	request.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("couldn't post url %s", postURL.String())
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("coudln't send metrics, status code: %d, response: %s", resp.StatusCode, body)
	}

	if err := resp.Body.Close(); err != nil {
		return errors.New("couldn't close response body")
	}

	return err
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
