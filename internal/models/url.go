package models

import (
	"net/http"

	"go-metricscol/internal/server/apierror"

	"github.com/go-chi/chi"
)

type GetURLData struct {
	MetricName string
	MetricType MetricType
}

type PostURLData struct {
	GetURLData
	MetricValue string
}

func ParsePostURLData(r *http.Request) (*PostURLData, error) {
	getData, apiError := ParseGetURLData(r)
	if apiError != nil {
		return nil, apiError
	}

	postData := PostURLData{
		GetURLData: *getData,
	}

	value := chi.URLParam(r, "value")

	if len(value) == 0 {
		return nil, apierror.EmptyArguments
	}
	postData.MetricValue = value

	return &postData, nil
}

func ParseGetURLData(r *http.Request) (*GetURLData, error) {
	urlData := GetURLData{}
	switch chi.URLParam(r, "type") {
	case "gauge":
		urlData.MetricType = Gauge
	case "counter":
		urlData.MetricType = Counter
	default:
		return nil, apierror.UnknownMetricType
	}

	name := chi.URLParam(r, "name")

	if len(name) == 0 {
		return nil, apierror.EmptyArguments
	}
	urlData.MetricName = name

	return &urlData, nil
}
