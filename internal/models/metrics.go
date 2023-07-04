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

func (m MetricsMap) Get(name string, valueType MetricType) (*Metric, apierror.APIError) {
	metric, ok := m[getKey(name, valueType)]
	if !ok {
		return nil, apierror.NotFound
	}

	return &metric, apierror.NoError
}

func (m MetricsMap) Update(name string, valueType MetricType, value interface{}) apierror.APIError {
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

		m[getKey(name, valueType)] = Metric{Name: name, MType: GaugeType, Value: utils.Ptr(floatValue)}
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
		var prevVal int64
		if prevMetric == nil {
			prevVal = 0
		} else {
			prevVal = *prevMetric.Delta
		}

		m[getKey(name, valueType)] = Metric{Name: name, MType: CounterType, Delta: utils.Ptr(prevVal + intValue)}
	default:
		return apierror.UnknownMetricType
	}

	return apierror.NoError
}

func (m MetricsMap) ResetPollCount() {
	m[getKey("PollCount", CounterType)] = Metric{Name: "PollCount", MType: CounterType, Delta: utils.Ptr(int64(0))}
}
