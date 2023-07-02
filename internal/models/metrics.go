package models

import (
	"errors"
	"fmt"
	"go-metricscol/internal/server/apierror"
	"log"
	"net/http"
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

func (m Metrics) Get(name string, valueType MetricType) (Metric, apierror.APIError) {
	metric, ok := m[getKey(name, valueType)]
	if !ok {
		return nil, apierror.NotFound
	}

	return metric, apierror.NoError
}

func (m Metrics) Update(name string, valueType MetricType, value interface{}) apierror.APIError {
	switch valueType {
	case GaugeType:
		if v, ok := value.(float64); ok {
			m[getKey(name, valueType)] = Gauge{Name: name, Value: v}
		} else {
			return apierror.InvalidValue
		}
	case CounterType:
		if v, ok := value.(int64); ok {
			prevMetric, _ := m.Get(name, CounterType)
			prevVal, _ := prevMetric.(Counter) // TODO: Насколько это безопасно?

			m[getKey(name, valueType)] = Counter{Name: name, Value: prevVal.Value + v}
		} else {
			return apierror.InvalidValue
		}
	default:
		return apierror.UnknownMetricType
	}

	return apierror.NoError
}

func (m Metrics) ResetPollCount() {
	m[getKey("PollCount", CounterType)] = Counter{Name: "PollCount", Value: 0}
}

func (m Metrics) SendToServer(addr string) error {
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
