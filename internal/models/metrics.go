package models

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"
	"log"
	"strings"
)

var key string

type Metrics struct {
	Collection map[string]Metric
	config     *Config
}

func init() {
	flag.StringVar(&key, "k", "", "Key to encrypt metrics")
}

func parseConfig() *Config {
	flag.Parse()
	config := NewConfig(key)

	if err := env.Parse(config); err != nil {
		log.Fatalf("Couldn't parse config with error: %s", err)
	}
	return config
}

func NewMetrics() Metrics {
	return Metrics{Collection: map[string]Metric{}, config: parseConfig()}
}

func getKey(name string, valueType MetricType) string {
	key := strings.Builder{}
	key.WriteString(name)
	key.WriteByte(':')
	key.WriteString(valueType.String())

	return key.String()
}

func (m Metrics) Get(name string, valueType MetricType) (*Metric, error) {
	metric, ok := m.Collection[getKey(name, valueType)]
	if !ok {
		return nil, apierror.NotFound
	}

	return &metric, nil
}

func (m Metrics) Update(name string, valueType MetricType, value interface{}) error {
	if valueType != Gauge && valueType != Counter {
		return apierror.UnknownMetricType
	}

	metricKey := getKey(name, valueType)
	var metric Metric

	switch valueType {
	case Gauge:
		var floatValue float64
		switch v := value.(type) {
		case float32:
			floatValue = float64(v)
		case float64:
			floatValue = v
		case int:
			floatValue = float64(v)
		case int8:
			floatValue = float64(v)
		case int16:
			floatValue = float64(v)
		case int32:
			floatValue = float64(v)
		case int64:
			floatValue = float64(v)
		default:
			return apierror.InvalidValue
		}

		metric = Metric{Name: name, MType: Gauge, Value: utils.Ptr(floatValue)}
	case Counter:
		var intValue int64
		switch v := value.(type) {
		case int:
			intValue = int64(v)
		case int8:
			intValue = int64(v)
		case int16:
			intValue = int64(v)
		case int32:
			intValue = int64(v)
		case int64:
			intValue = v
		default:
			return apierror.InvalidValue
		}

		prevMetric, _ := m.Get(name, Counter)
		var prevVal int64
		if prevMetric == nil {
			prevVal = 0
		} else {
			prevVal = *prevMetric.Delta
		}

		metric = Metric{Name: name, MType: Counter, Delta: utils.Ptr(prevVal + intValue)}
	default:
		return apierror.UnknownMetricType
	}

	metric.SetHashValue(m.config.CryptoKey)
	m.Collection[metricKey] = metric
	return nil
}

func (m Metrics) UpdateWithStruct(metric *Metric) error {
	if metric.MType != Gauge && metric.MType != Counter {
		return apierror.UnknownMetricType
	}

	metric.SetHashValue(m.config.CryptoKey)

	switch metric.MType {
	case Gauge:
		if metric.Value == nil || metric.Delta != nil {
			return apierror.InvalidValue
		}

		m.Collection[getKey(metric.Name, metric.MType)] = *metric
	case Counter:
		if metric.Delta == nil || metric.Value != nil {
			return apierror.InvalidValue
		}

		prevMetric, _ := m.Get(metric.Name, Counter)
		var prevVal, currentVal int64
		if prevMetric == nil {
			prevVal = 0
		} else {
			prevVal = *prevMetric.Delta
		}

		if metric.Delta == nil {
			currentVal = 0
		} else {
			currentVal = *metric.Delta
		}

		m.Collection[getKey(metric.Name, metric.MType)] = Metric{Name: metric.Name, MType: Counter, Delta: utils.Ptr(prevVal + currentVal), Hash: metric.Hash}
	default:
		return apierror.UnknownMetricType
	}

	return nil
}

func (m Metrics) ResetPollCount() {
	m.Collection[getKey("PollCount", Counter)] = Metric{Name: "PollCount", MType: Counter, Delta: utils.Ptr(int64(0))}
}
