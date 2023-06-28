package models

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

//TODO: Точно ли тут float64?

type MetricType int

const (
	Gauge   MetricType = iota // float64
	Counter                   //int64
)

func (m MetricType) String() string {
	switch m {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	}
	return ""
}

type Metric struct {
	valueType MetricType
	value     uint64
}

func (m Metric) StringValue() string {
	switch m.valueType {
	case Gauge:
		return fmt.Sprintf("%v", math.Float64frombits(m.value))
	case Counter:
		return strconv.FormatInt(int64(m.value), 10)
	}

	return ""
}

func NewMetric(valueType MetricType) Metric {
	return Metric{valueType: valueType}
}

type Metrics map[string]Metric

func (m Metrics) SendToServer(addr string) {
	for name, metric := range m {
		postURL := fmt.Sprintf("%s/update/%s/%s/%s", addr, metric.valueType.String(), name, metric.StringValue())
		log.Println(postURL)
		resp, err := http.Post(postURL, "text/plain", nil)

		if err := resp.Body.Close(); err != nil {
			log.Fatalf("Couldn't close response body")
		}

		if err != nil {
			log.Fatalf("Couldn't send metric %s to server", name)
		}
	}
}

func (m Metrics) UpdateGauge(key string, value float64) {
	metric, ok := m[key]
	if !ok {
		log.Panic("No such key in metrics map")
	}

	if metric.valueType != Gauge {
		log.Panic("Types are mismatched in metrics map")
	}

	metric.value = math.Float64bits(value)
	m[key] = metric
}

func (m Metrics) UpdateCounter(key string) {
	metric, ok := m[key]
	if !ok {
		log.Panic("No such key in metrics map")
	}

	if metric.valueType != Counter {
		log.Panic("Types are mismatched in metrics map")
	}

	metric.value++
	m[key] = metric
}

func NewMetrics() Metrics {
	metrics := make(Metrics, 29)
	metrics["Alloc"] = NewMetric(Gauge)
	metrics["BuckHashSys"] = NewMetric(Gauge)
	metrics["Frees"] = NewMetric(Gauge)
	metrics["GCCPUFraction"] = NewMetric(Gauge)
	metrics["GCSys"] = NewMetric(Gauge)
	metrics["HeapAlloc"] = NewMetric(Gauge)
	metrics["HeapIdle"] = NewMetric(Gauge)
	metrics["HeapInuse"] = NewMetric(Gauge)
	metrics["HeapObjects"] = NewMetric(Gauge)
	metrics["HeapReleased"] = NewMetric(Gauge)
	metrics["HeapSys"] = NewMetric(Gauge)
	metrics["LastGC"] = NewMetric(Gauge)
	metrics["Lookups"] = NewMetric(Gauge)
	metrics["MCacheInuse"] = NewMetric(Gauge)
	metrics["MCacheSys"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["MSpanSys"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["Mallocs"] = NewMetric(Gauge)
	metrics["NextGC"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)

	metrics["Alloc"] = NewMetric(Gauge)
	metrics["BuckHashSys"] = NewMetric(Gauge)
	metrics["Frees"] = NewMetric(Gauge)
	metrics["GCCPUFraction"] = NewMetric(Gauge)
	metrics["GCSys"] = NewMetric(Gauge)
	metrics["HeapAlloc"] = NewMetric(Gauge)
	metrics["HeapIdle"] = NewMetric(Gauge)
	metrics["HeapInuse"] = NewMetric(Gauge)
	metrics["HeapObjects"] = NewMetric(Gauge)
	metrics["HeapReleased"] = NewMetric(Gauge)
	metrics["HeapSys"] = NewMetric(Gauge)
	metrics["LastGC"] = NewMetric(Gauge)
	metrics["Lookups"] = NewMetric(Gauge)
	metrics["MCacheInuse"] = NewMetric(Gauge)
	metrics["MCacheSys"] = NewMetric(Gauge)
	metrics["MSpanInuse"] = NewMetric(Gauge)
	metrics["MSpanSys"] = NewMetric(Gauge)
	metrics["Mallocs"] = NewMetric(Gauge)
	metrics["NextGC"] = NewMetric(Gauge)
	metrics["NumForcedGC"] = NewMetric(Gauge)
	metrics["NumGC"] = NewMetric(Gauge)
	metrics["OtherSys"] = NewMetric(Gauge)
	metrics["PauseTotalNs"] = NewMetric(Gauge)
	metrics["StackInuse"] = NewMetric(Gauge)
	metrics["StackSys"] = NewMetric(Gauge)
	metrics["Sys"] = NewMetric(Gauge)
	metrics["TotalAlloc"] = NewMetric(Gauge)
	metrics["RandomValue"] = NewMetric(Gauge)
	metrics["PollCount"] = NewMetric(Counter)

	return metrics
}
