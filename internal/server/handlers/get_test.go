package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers_Get(t *testing.T) {
	h := Handlers{
		Storage: memory.NewMemStorage(),
	}

	require.NoError(t, h.Storage.Update("Alloc", models.GaugeType, "123.4"))
	require.NoError(t, h.Storage.Update("MemoryInUse", models.GaugeType, "593"))
	require.NoError(t, h.Storage.Update("PollCount", models.CounterType, "1"))

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
				metricType: models.GaugeType.String(),
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
				metricType: models.GaugeType.String(),
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

	require.NoError(t, h.Storage.Update("Alloc", models.GaugeType, "123.4"))
	require.NoError(t, h.Storage.Update("MemoryInUse", models.GaugeType, "593"))
	require.NoError(t, h.Storage.Update("PollCount", models.CounterType, "1"))

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
