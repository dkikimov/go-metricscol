package models

import (
	"errors"
	"fmt"
	"go-metricscol/internal/server/apiError"
	"log"
	"math"
	"net/http"
	"strconv"
)

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

func (m Metric) ValueType() MetricType {
	return m.valueType
}

func NewMetric(valueType MetricType) Metric {
	return Metric{valueType: valueType}
}

type Metrics map[string]Metric

func (m Metrics) SendToServer(addr string) error {
	for name, metric := range m {
		postURL := fmt.Sprintf("%s/update/%s/%s/%s", addr, metric.valueType.String(), name, metric.StringValue())
		log.Println(postURL)
		resp, err := http.Post(postURL, "text/plain", nil)

		if err != nil {
			return errors.New(fmt.Sprintf("couldn't post url %s", postURL))
		}

		if err := resp.Body.Close(); err != nil {
			return errors.New("couldn't close response body")
		}

	}
	return nil
}

func (m Metrics) UpdateGauge(key string, value float64) apiError.APIError {
	metric, ok := m[key]
	if !ok {
		metric = NewMetric(Gauge)
	}

	if metric.valueType != Gauge {
		return apiError.TypeMismatch
	}

	metric.value = math.Float64bits(value)
	m[key] = metric

	return apiError.NoError
}

func (m Metrics) UpdateCounter(key string, value int64) apiError.APIError {
	metric, ok := m[key]
	if !ok {
		metric = NewMetric(Counter)
	}

	if metric.valueType != Counter {
		return apiError.TypeMismatch
	}

	metric.value += uint64(value)
	m[key] = metric
	return apiError.NoError
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
