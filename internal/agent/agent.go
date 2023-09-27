package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"golang.org/x/sync/errgroup"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
)

// SendMetricsToServer sends metrics stored is memory.Metrics to the address given in agent.Config.
// Rate limit defined in config is not exceeded.
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

// UpdateMetrics gets all metrics from runtime.MemStats and writes them to memory.Metrics.
func UpdateMetrics(metrics *memory.Metrics) error {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	if err := metrics.Update("Alloc", models.Gauge, float64(stats.Alloc)); err != nil {
		return fmt.Errorf("couldn't collect Alloc: %s", err)
	}
	if err := metrics.Update("BuckHashSys", models.Gauge, float64(stats.BuckHashSys)); err != nil {
		return fmt.Errorf("couldn't collect BuckHashSys: %s", err)
	}
	if err := metrics.Update("Frees", models.Gauge, float64(stats.Frees)); err != nil {
		return fmt.Errorf("couldn't collect Frees: %s", err)
	}
	if err := metrics.Update("GCCPUFraction", models.Gauge, stats.GCCPUFraction); err != nil {
		return fmt.Errorf("couldn't collect GCCPUFraction: %s", err)
	}
	if err := metrics.Update("GCSys", models.Gauge, float64(stats.GCSys)); err != nil {
		return fmt.Errorf("couldn't collect GCSys: %s", err)
	}
	if err := metrics.Update("HeapAlloc", models.Gauge, float64(stats.HeapAlloc)); err != nil {
		return fmt.Errorf("couldn't collect HeapAlloc: %s", err)
	}
	if err := metrics.Update("HeapIdle", models.Gauge, float64(stats.HeapIdle)); err != nil {
		return fmt.Errorf("couldn't collect HeapIdle: %s", err)
	}
	if err := metrics.Update("HeapInuse", models.Gauge, float64(stats.HeapInuse)); err != nil {
		return fmt.Errorf("couldn't collect HeapInuse: %s", err)
	}
	if err := metrics.Update("HeapObjects", models.Gauge, float64(stats.HeapObjects)); err != nil {
		return fmt.Errorf("couldn't collect HeapObjects: %s", err)
	}
	if err := metrics.Update("HeapReleased", models.Gauge, float64(stats.HeapReleased)); err != nil {
		return fmt.Errorf("couldn't collect HeapReleased: %s", err)
	}
	if err := metrics.Update("HeapSys", models.Gauge, float64(stats.HeapSys)); err != nil {
		return fmt.Errorf("couldn't collect HeapSys: %s", err)
	}
	if err := metrics.Update("LastGC", models.Gauge, float64(stats.LastGC)); err != nil {
		return fmt.Errorf("couldn't collect LastGC: %s", err)
	}
	if err := metrics.Update("Lookups", models.Gauge, float64(stats.Lookups)); err != nil {
		return fmt.Errorf("couldn't collect Lookups: %s", err)
	}
	if err := metrics.Update("MCacheInuse", models.Gauge, float64(stats.MCacheInuse)); err != nil {
		return fmt.Errorf("couldn't collect MCacheInuse: %s", err)
	}
	if err := metrics.Update("MCacheSys", models.Gauge, float64(stats.MCacheSys)); err != nil {
		return fmt.Errorf("couldn't collect MCacheSys: %s", err)
	}
	if err := metrics.Update("MSpanInuse", models.Gauge, float64(stats.MSpanInuse)); err != nil {
		return fmt.Errorf("couldn't collect MSpanInuse: %s", err)
	}
	if err := metrics.Update("MSpanSys", models.Gauge, float64(stats.MSpanSys)); err != nil {
		return fmt.Errorf("couldn't collect MSpanSys: %s", err)
	}
	if err := metrics.Update("Mallocs", models.Gauge, float64(stats.Mallocs)); err != nil {
		return fmt.Errorf("couldn't collect Mallocs: %s", err)
	}
	if err := metrics.Update("NextGC", models.Gauge, float64(stats.NextGC)); err != nil {
		return fmt.Errorf("couldn't collect NextGC: %s", err)
	}
	if err := metrics.Update("NumForcedGC", models.Gauge, float64(stats.NumForcedGC)); err != nil {
		return fmt.Errorf("couldn't collect NumForcedGC: %s", err)
	}
	if err := metrics.Update("NumGC", models.Gauge, float64(stats.NumGC)); err != nil {
		return fmt.Errorf("couldn't collect NumGC: %s", err)
	}
	if err := metrics.Update("OtherSys", models.Gauge, float64(stats.OtherSys)); err != nil {
		return fmt.Errorf("couldn't collect OtherSys: %s", err)
	}
	if err := metrics.Update("PauseTotalNs", models.Gauge, float64(stats.PauseTotalNs)); err != nil {
		return fmt.Errorf("couldn't collect PauseTotalNs: %s", err)
	}
	if err := metrics.Update("StackInuse", models.Gauge, float64(stats.StackInuse)); err != nil {
		return fmt.Errorf("couldn't collect StackInuse: %s", err)
	}
	if err := metrics.Update("StackSys", models.Gauge, float64(stats.StackSys)); err != nil {
		return fmt.Errorf("couldn't collect StackSys: %s", err)
	}
	if err := metrics.Update("Sys", models.Gauge, float64(stats.Sys)); err != nil {
		return fmt.Errorf("couldn't collect Sys: %s", err)
	}
	if err := metrics.Update("TotalAlloc", models.Gauge, float64(stats.TotalAlloc)); err != nil {
		return fmt.Errorf("couldn't collect TotalAlloc: %s", err)
	}
	if err := metrics.Update("RandomValue", models.Gauge, rand.Float64()); err != nil {
		return fmt.Errorf("couldn't collect RandomValue: %s", err)
	}
	if err := metrics.Update("PollCount", models.Counter, 1); err != nil {
		return fmt.Errorf("couldn't collect PollCount: %s", err)
	}

	return nil
}

// CollectAdditionalMetrics writes memory and CPU usage metrics to the memory.Metrics.
func CollectAdditionalMetrics(metrics *memory.Metrics) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("couldn't collect memory with error: %s", err)
	}

	if err := metrics.Update("TotalMemory", models.Gauge, float64(v.Total)); err != nil {
		return fmt.Errorf("couldn't collect TotalMemory: %s", err)
	}

	if err := metrics.Update("FreeMemory", models.Gauge, float64(v.Free)); err != nil {
		return fmt.Errorf("couldn't collect FreeMemory: %s", err)
	}

	coresPercent, err := cpu.Percent(time.Second, true)
	if err != nil {
		return fmt.Errorf("couldn't collect cpu utilization with error: %s", err)
	}

	for i, core := range coresPercent {
		num := i + 1
		if err := metrics.Update(fmt.Sprintf("CPUutilization%d", num), models.Gauge, core); err != nil {
			return fmt.Errorf("couldn't collect CPUutilization%d: %s", num, err)
		}
	}

	return nil
}
