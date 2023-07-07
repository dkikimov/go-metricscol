package models

import (
	"go-metricscol/internal/server/apierror"
	"strings"
)

type Metrics map[string]Metric

func getKey(name string, valueType MetricType) string {
	key := strings.Builder{}
	key.WriteString(name)
	key.WriteByte(':')
	key.WriteString(valueType.String())

	return key.String()
}

func (m Metrics) Get(name string, valueType MetricType) (Metric, error) {
	metric, ok := m[getKey(name, valueType)]
	if !ok {
		return nil, apierror.NotFound
	}

	return metric, nil
}

func (m Metrics) Update(name string, valueType MetricType, value interface{}) error {
	if valueType != GaugeType && valueType != CounterType {
		return apierror.UnknownMetricType
	}

	switch valueType {
	case GaugeType:
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

		m[getKey(name, valueType)] = Gauge{Name: name, Value: floatValue}
	case CounterType:
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
		prevMetric, _ := m.Get(name, CounterType)
		prevVal, _ := prevMetric.(Counter)

		m[getKey(name, valueType)] = Counter{Name: name, Value: prevVal.Value + intValue}
	default:
		return apierror.UnknownMetricType
	}

	return nil
}

func (m Metrics) ResetPollCount() {
	m[getKey("PollCount", CounterType)] = Counter{Name: "PollCount", Value: 0}
}
