package grpcpackage

import (
	"strconv"

	"go-metricscol/internal/models"
	"go-metricscol/internal/proto"
	"go-metricscol/internal/server/apierror"
)

func parseMetricFromRequest(metric *proto.Metric) (*models.Metric, error) {
	var resultMetric models.Metric

	resultMetric.Name = metric.Name
	switch metric.Type {
	case proto.MetricType_GAUGE:
		resultMetric.MType = models.Gauge

		floatVal, err := strconv.ParseFloat(metric.Value, 64)
		if err != nil {
			return nil, apierror.NumberParse
		}

		resultMetric.Value = &floatVal
	case proto.MetricType_COUNTER:
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

func parseTypeFromRequest(metricType proto.MetricType) (models.MetricType, error) {
	switch metricType {
	case proto.MetricType_GAUGE:
		return models.Gauge, nil
	case proto.MetricType_COUNTER:
		return models.Counter, nil
	default:
		return "", apierror.UnknownMetricType
	}
}
