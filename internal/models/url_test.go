package models

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"go-metricscol/internal/server/apierror"
	"net/http"
	"testing"
)

type urlParams struct {
	metricName  string
	metricType  MetricType
	metricValue string
}

func getRequest(method string, url string, params urlParams) *http.Request {
	request, err := http.NewRequest(method, fmt.Sprintf("/%s/%s/%s/%s", url, params.metricType, params.metricName, params.metricValue), nil)

	if err != nil {
		panic(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", params.metricType.String())
	rctx.URLParams.Add("name", params.metricName)
	rctx.URLParams.Add("value", params.metricValue)

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
	return request
}

func TestParseGetURLData(t *testing.T) {

	tests := []struct {
		name    string
		args    urlParams
		want    *GetURLData
		wantErr error
	}{
		{
			name: "Get metric gauge",
			args: urlParams{
				metricName: "Alloc",
				metricType: GaugeType,
			},
			want: &GetURLData{
				MetricName: "Alloc",
				MetricType: GaugeType,
			},
		},
		{
			name: "Get metric counter",
			args: urlParams{
				metricName: "PollCount",
				metricType: CounterType,
			},
			want: &GetURLData{
				MetricName: "PollCount",
				MetricType: CounterType,
			},
		},
		{
			name: "Empty name",
			args: urlParams{
				metricName: "",
				metricType: GaugeType,
			},
			want:    nil,
			wantErr: apierror.EmptyArguments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGetURLData(getRequest(http.MethodGet, "value", tt.args))
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsePostURLData(t *testing.T) {
	tests := []struct {
		name    string
		args    urlParams
		want    *PostURLData
		wantErr error
	}{
		{
			name: "Post metric gauge",
			args: urlParams{
				metricName:  "Alloc",
				metricType:  GaugeType,
				metricValue: "13",
			},
			want: &PostURLData{
				GetURLData: GetURLData{
					MetricName: "Alloc",
					MetricType: GaugeType,
				},
				MetricValue: "13",
			},
		},
		{
			name: "Post metric counter",
			args: urlParams{
				metricName:  "PollCount",
				metricType:  CounterType,
				metricValue: "1",
			},
			want: &PostURLData{
				GetURLData: GetURLData{
					MetricName: "PollCount",
					MetricType: CounterType,
				},
				MetricValue: "1",
			},
		},
		{
			name: "Empty name",
			args: urlParams{
				metricName:  "",
				metricType:  GaugeType,
				metricValue: "1",
			},
			want:    nil,
			wantErr: apierror.EmptyArguments,
		},
		{
			name: "Empty value",
			args: urlParams{
				metricName:  "Alloc",
				metricType:  GaugeType,
				metricValue: "",
			},
			want:    nil,
			wantErr: apierror.EmptyArguments,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePostURLData(getRequest(http.MethodPost, "update", tt.args))
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
