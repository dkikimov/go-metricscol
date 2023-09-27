package models

import (
	"net/http"

	"github.com/go-chi/chi"

	"go-metricscol/internal/server/apierror"
)

// GetURLData describes the parameters passed to the URL in the GET request.
type GetURLData struct {
	MetricName string
	MetricType MetricType
}

// PostURLData describes the parameters passed to the URL in the POST request.
type PostURLData struct {
	GetURLData
	MetricValue string
}

// ParsePostURLData parses parameters passed to the URL and saves it into PostURLData struct.
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

// ParseGetURLData parses parameters passed to the URL and saves it into GetURLData struct.
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
