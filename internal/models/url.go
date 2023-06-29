package models

import (
	"go-metricscol/internal/server/apiError"
	"strings"
)

type URLData struct {
	MetricName  string
	MetricValue string
	MetricType  MetricType
}

func ParseURLData(url string) (*URLData, apiError.APIError) {
	splitURL := strings.Split(url, "/")[2:]
	if len(splitURL) < 3 {
		return nil, apiError.NotEnoughArguments
	}

	urlData := URLData{}
	data := splitURL[len(splitURL)-3:]
	switch data[0] {
	case "gauge":
		urlData.MetricType = Gauge
	case "counter":
		urlData.MetricType = Counter
	default:
		return nil, apiError.UnknownMetricType
	}

	if len(data[1]) == 0 || len(data[2]) == 0 {
		return nil, apiError.EmptyArguments
	}
	urlData.MetricName = data[1]
	urlData.MetricValue = data[2]

	return &urlData, apiError.NoError
}
