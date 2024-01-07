package proto

import (
	"strconv"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
)

func ParseMetricFromRequest(metric *Metric) (*models.Metric, error) {
	var resultMetric models.Metric

	resultMetric.Name = metric.Name
	resultMetric.Hash = metric.Hash
	switch metric.Type {
	case MetricType_GAUGE:
		resultMetric.MType = models.Gauge

		floatVal, err := strconv.ParseFloat(metric.Value, 64)
		if err != nil {
			return nil, apierror.NumberParse
		}

		resultMetric.Value = &floatVal
	case MetricType_COUNTER:
		resultMetric.MType = models.Counter

		intVal, err := strconv.ParseInt(metric.Value, 10, 64)
		if err != nil {
			return nil, apierror.NumberParse
		}

		resultMetric.Delta = &intVal
	default:
		return nil, apierror.UnknownMetricType
	}

	return &resultMetric, nil
}

func ParseTypeFromRequest(metricType MetricType) (models.MetricType, error) {
	switch metricType {
	case MetricType_GAUGE:
		return models.Gauge, nil
	case MetricType_COUNTER:
		return models.Counter, nil
	default:
		return "", apierror.UnknownMetricType
	}
}
