package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/utils"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlers_Update(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "Without name and value counter",
			args: args{
				metricType: models.Counter.String(),
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Without value counter",
			args: args{
				metricType: models.Counter.String(),
				metricName: "Alloc",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Without name and value gauge",
			args: args{
				metricType: models.Gauge.String(),
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Without value gauge",
			args: args{
				metricType: models.Gauge.String(),
				metricName: "Alloc",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Wrong value type",
			args: args{
				metricType:  models.Gauge.String(),
				metricName:  "PollCount",
				metricValue: "hello",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Metric name and type mismatch counter",
			args: args{
				metricType:  models.Counter.String(),
				metricName:  "Alloc",
				metricValue: "1.23",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(
				memory.NewMemStorage(),
				nil,
				NewConfig(""),
			)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/value/%s/%s/%s", tt.args.metricType, tt.args.metricName, tt.args.metricValue), nil)
			if err != nil {
				t.Fatal(err)
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", tt.args.metricType)
			rctx.URLParams.Add("name", tt.args.metricName)
			rctx.URLParams.Add("value", tt.args.metricValue)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Update)

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatusCode, rr.Code)
		})
	}
}

func TestHandlers_UpdateJSON(t *testing.T) {
	type want struct {
		Body       models.Metric
		StatusCode int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "Update gauge",
			body: `{"id": "Alloc", "type": "gauge", "value": 13.1}`,
			want: want{
				Body: models.Metric{
					Name:  "Alloc",
					MType: models.Gauge,
					Value: utils.Ptr(13.1),
				},
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Update counter",
			body: `{"id": "PollCount", "type": "counter", "delta": 13}`,
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
			name: "Update unknown type",
			body: `{"id": "Alloc", "type": "unknown", "value": 13.1}`,
			want: want{
				StatusCode: http.StatusNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		storage := memory.NewMemStorage()
		h := NewHandlers(
			storage,
			nil,
			NewConfig(""),
		)

		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/update/", bytes.NewReader([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.UpdateJSON)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.want.StatusCode, rr.Code)
			if rr.Code == http.StatusOK {
				got, _ := storage.GetAll(context.Background())
				require.Equal(t, 1, len(got))

				assert.True(t, reflect.DeepEqual(tt.want.Body, got[0]))
			}
		})
	}
}

func BenchmarkHandlers_Update_MemStorage(b *testing.B) {
	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig("hash"),
	)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/value/%s/%s/%s", "Alloc", models.Gauge, "121.14"), nil)
	require.NoError(b, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", models.Gauge.String())
	rctx.URLParams.Add("name", "Alloc")
	rctx.URLParams.Add("value", "121.14")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Update)

	// Выключить логи для handler'а
	log.SetOutput(io.Discard)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
		assert.Equal(b, 200, rr.Code)
	}
}

func BenchmarkHandlers_UpdateJSON_MemStorage(b *testing.B) {
	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig("hash"),
	)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateJSON)

	// Выключить логи для handler'а
	log.SetOutput(io.Discard)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewReader([]byte(`{"id": "Alloc", "type": "gauge", "value": 13.1}`)))
		handler.ServeHTTP(rr, req)
		assert.Equal(b, 200, rr.Code)
	}
}
