package models

import (
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"
	"strings"
)

type MetricsMap map[string]Metric

func getKey(name string, valueType MetricType) string {
	key := strings.Builder{}
	key.WriteString(name)
	key.WriteByte(':')
	key.WriteString(valueType.String())

	return key.String()
}

func (m MetricsMap) Get(name string, valueType MetricType) (*Metric, error) {
	metric, ok := m[getKey(name, valueType)]
	if !ok {
		return nil, apierror.NotFound
	}

	return &metric, nil
}

func (m MetricsMap) Update(name string, valueType MetricType, value interface{}) error {
	if valueType != Gauge && valueType != Counter {
		return apierror.UnknownMetricType
	}

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

		m[getKey(name, valueType)] = Metric{Name: name, MType: Gauge, Value: utils.Ptr(floatValue)}
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

		m[getKey(name, valueType)] = Metric{Name: name, MType: Counter, Delta: utils.Ptr(prevVal + intValue)}
	default:
		return apierror.UnknownMetricType
	}

	return nil
}

func (m MetricsMap) UpdateWithStruct(metric *Metric) error {
	if metric.MType != Gauge && metric.MType != Counter {
		return apierror.UnknownMetricType
	}
	if len(metric.Name) == 0 {
		return apierror.InvalidValue
	}

	switch metric.MType {
	case Gauge:
		if metric.Value == nil {
			return apierror.InvalidValue
		}

		m[getKey(metric.Name, metric.MType)] = *metric
	case Counter:
		if metric.Delta == nil {
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

		m[getKey(metric.Name, metric.MType)] = Metric{Name: metric.Name, MType: Counter, Delta: utils.Ptr(prevVal + currentVal)}
	default:
		return apierror.UnknownMetricType
	}

	return nil
}

func (m MetricsMap) ResetPollCount() {
	m[getKey("PollCount", Counter)] = Metric{Name: "PollCount", MType: Counter, Delta: utils.Ptr(int64(0))}
}
