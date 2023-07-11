package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/utils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandlers_Get(t *testing.T) {
	h := Handlers{
		Storage: memory.NewMemStorage(),
	}

	require.NoError(t, h.Storage.Update("Alloc", models.Gauge, "123.4"))
	require.NoError(t, h.Storage.Update("MemoryInUse", models.Gauge, "593"))
	require.NoError(t, h.Storage.Update("PollCount", models.Counter, "1"))

	type want struct {
		StatusCode int
		Body       string
	}
	type args struct {
		metricType string
		metricName string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get Alloc value",
			args: args{
				metricType: models.Gauge.String(),
				metricName: "Alloc",
			},
			want: want{
				StatusCode: http.StatusOK,
				Body:       "123.4",
			},
		},
		{
			name: "Unknown metric",
			args: args{
				metricType: models.Gauge.String(),
				metricName: "NewMetric",
			},
			want: want{
				StatusCode: http.StatusNotFound,
				Body:       "",
			},
		},
		{
			name: "Unknown metric type",
			args: args{
				metricType: "h",
				metricName: "Alloc",
			},
			want: want{
				StatusCode: http.StatusNotImplemented,
				Body:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/value/%s/%s", tt.args.metricType, tt.args.metricName), nil)
			if err != nil {
				t.Fatal(err)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.args.metricType)
			rctx.URLParams.Add("name", tt.args.metricName)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Get)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.want.StatusCode, rr.Code)
			assert.Equal(t, tt.want.Body, rr.Body.String())
		})
	}
}

func TestHandlers_GetAll(t *testing.T) {
	h := Handlers{
		Storage: memory.NewMemStorage(),
	}

	require.NoError(t, h.Storage.Update("Alloc", models.Gauge, "123.4"))
	require.NoError(t, h.Storage.Update("MemoryInUse", models.Gauge, "593"))
	require.NoError(t, h.Storage.Update("PollCount", models.Counter, "1"))

	type want struct {
		StatusCode int
		Body       string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get all values",
			args: args{
				url: "/",
			},
			want: want{
				StatusCode: http.StatusOK,
				Body: "Key: Alloc, value: 123.4, type: gauge \n" +
					"Key: MemoryInUse, value: 593, type: gauge \n" +
					"Key: PollCount, value: 1, type: counter \n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.args.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetAll)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.want.StatusCode, rr.Code)
			require.Equal(t, tt.want.Body, rr.Body.String())
		})
	}
}

func TestHandlers_GetJSON(t *testing.T) {
	storage := memory.NewMemStorage()
	require.NoError(t, storage.Update("Alloc", "gauge", "12.1"))
	require.NoError(t, storage.Update("PollCount", "counter", "13"))

	h := Handlers{Storage: storage}
	type want struct {
		Body       models.Metric
		StatusCode int
	}

	tests := []struct {
		name string

		//TODO: Как сделать предподчительнее: писать сырой json, или маршалить из структуры?
		body models.Metric
		want want
	}{
		{
			name: "Get gauge",
			body: models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
			},
			want: want{
				Body: models.Metric{
					Name:  "Alloc",
					MType: models.Gauge,
					Value: utils.Ptr(12.1),
				},
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Get counter",
			body: models.Metric{
				Name:  "PollCount",
				MType: models.Counter,
			},
			want: want{
				Body: models.Metric{
					Name:  "PollCount",
					MType: models.Counter,
					Delta: utils.Ptr(int64(13)),
				},
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Get unknown metric",
			body: models.Metric{
				Name:  "H",
				MType: models.Counter,
			},
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Get wrong type",
			body: models.Metric{
				Name:  "PollCount",
				MType: "unknown",
			},
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Get gauge with counter type",
			body: models.Metric{
				Name:  "Alloc",
				MType: models.Counter,
			},
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Get counter with gauge type",
			body: models.Metric{
				Name:  "PollCount",
				MType: models.Gauge,
			},
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricJson, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/value/", bytes.NewReader(metricJson))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetJSON)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.want.StatusCode, rr.Code)
			if rr.Code == http.StatusOK {
				var gotMetric models.Metric
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &gotMetric))
				assert.True(t, reflect.DeepEqual(tt.want.Body, gotMetric))
			}
		})
	}
}
