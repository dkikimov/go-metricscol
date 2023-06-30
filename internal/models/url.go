package models

import (
	"github.com/go-chi/chi"
	"go-metricscol/internal/server/apierror"
	"net/http"
)

type URLData struct {
	MetricName  string
	MetricValue string
	MetricType  MetricType
}

func ParseURLData(r *http.Request) (*URLData, apierror.APIError) {
	urlData := URLData{}
	switch chi.URLParam(r, "type") {
	case "gauge":
		urlData.MetricType = Gauge
	case "counter":
		urlData.MetricType = Counter
	default:
		return nil, apierror.UnknownMetricType
	}

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	if len(name) == 0 || len(value) == 0 {
		return nil, apierror.EmptyArguments
	}
	urlData.MetricName = name
	urlData.MetricValue = value

	return &urlData, apierror.NoError
}
