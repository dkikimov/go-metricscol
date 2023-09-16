package handlers

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestHandlers_Find(t *testing.T) {
	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig(""),
	)

	require.NoError(t, h.Storage.Update(context.Background(), "Alloc", models.Gauge, "123.4"))
	require.NoError(t, h.Storage.Update(context.Background(), "MemoryInUse", models.Gauge, "593"))
	require.NoError(t, h.Storage.Update(context.Background(), "PollCount", models.Counter, "1"))

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
			name: "Find Alloc value",
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
			handler := http.HandlerFunc(h.Find)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.want.StatusCode, rr.Code)
			assert.Equal(t, tt.want.Body, rr.Body.String())
		})
	}
}

func TestHandlers_GetAll(t *testing.T) {
	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig(""),
	)

	require.NoError(t, h.Storage.Update(context.Background(), "Alloc", models.Gauge, "123.4"))
	require.NoError(t, h.Storage.Update(context.Background(), "MemoryInUse", models.Gauge, "593"))
	require.NoError(t, h.Storage.Update(context.Background(), "PollCount", models.Counter, "1"))

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
			name: "Find all values",
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

// TODO: Добавить тесты для постгреса

func TestHandlers_GetAllWithHash(t *testing.T) {
	hashKey := "test"

	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig(hashKey),
	)

	alloc := models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(123.4),
	}

	memoryInUse := models.Metric{
		Name:  "MemoryInUse",
		MType: models.Gauge,
		Value: utils.Ptr(float64(593)),
	}

	pollCount := models.Metric{
		Name:  "PollCount",
		MType: models.Counter,
		Delta: utils.Ptr(int64(1)),
	}

	require.NoError(t, h.Storage.UpdateWithStruct(context.Background(), &alloc))
	require.NoError(t, h.Storage.UpdateWithStruct(context.Background(), &memoryInUse))
	require.NoError(t, h.Storage.UpdateWithStruct(context.Background(), &pollCount))

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
			name: "Find all values",
			args: args{
				url: "/",
			},
			want: want{
				StatusCode: http.StatusOK,
				Body: fmt.Sprintf("Key: Alloc, value: 123.4, type: gauge, hash: %s \n"+
					"Key: MemoryInUse, value: 593, type: gauge, hash: %s \n"+
					"Key: PollCount, value: 1, type: counter, hash: %s \n", alloc.HashValue(hashKey), memoryInUse.HashValue(hashKey), pollCount.HashValue(hashKey)),
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

func TestHandlers_FindJSON(t *testing.T) {
	storage := memory.NewMemStorage()
	require.NoError(t, storage.Update(context.Background(), "Alloc", "gauge", "12.1"))
	require.NoError(t, storage.Update(context.Background(), "PollCount", "counter", "13"))

	h := NewHandlers(
		storage,
		nil,
		NewConfig(""),
	)
	type want struct {
		Body       models.Metric
		StatusCode int
	}

	tests := []struct {
		name string

		// TODO: Как сделать предподчительнее: писать сырой json, или маршалить из структуры?
		body models.Metric
		want want
	}{
		{
			name: "Find gauge",
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
			name: "Find counter",
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
			name: "Find unknown metric",
			body: models.Metric{
				Name:  "H",
				MType: models.Counter,
			},
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Find wrong type",
			body: models.Metric{
				Name:  "PollCount",
				MType: "unknown",
			},
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Find gauge with counter type",
			body: models.Metric{
				Name:  "Alloc",
				MType: models.Counter,
			},
			want: want{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Find counter with gauge type",
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
			metricJSON, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/value/", bytes.NewReader(metricJSON))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.FindJSON)

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

func BenchmarkHandlers_Find_MemStorage(b *testing.B) {
	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig("hash"),
	)

	require.NoError(b, h.Storage.Update(context.Background(), "Alloc", models.Gauge, "123.4"))

	req, err := http.NewRequest("GET", fmt.Sprintf("/value/%s/%s", models.Gauge, "Alloc"), nil)
	require.NoError(b, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", models.Gauge.String())
	rctx.URLParams.Add("name", "Alloc")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Find)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
		assert.Equal(b, 200, rr.Code)
	}
}

func BenchmarkHandlers_FindJSON_MemStorage(b *testing.B) {
	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig("hash"),
	)

	metric := models.Metric{Name: "Alloc", MType: models.Gauge, Value: utils.Ptr(123.4)}
	require.NoError(b, h.Storage.UpdateWithStruct(context.Background(), &metric))

	metricJSON, err := json.Marshal(models.Metric{Name: "Alloc", MType: models.Gauge})
	require.NoError(b, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.FindJSON)

	log.SetOutput(io.Discard)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, "/value/", bytes.NewReader(metricJSON))
		require.NoError(b, err)
		handler.ServeHTTP(rr, req)
		assert.Equal(b, 200, rr.Code)
	}
}

func BenchmarkHandlers_FindAllWithHash_MemStorage(b *testing.B) {
	h := NewHandlers(
		memory.NewMemStorage(),
		nil,
		NewConfig("hash"),
	)

	require.NoError(b, h.Storage.Update(context.Background(), "Alloc", models.Gauge, "123.4"))
	require.NoError(b, h.Storage.Update(context.Background(), "Mem", models.Gauge, "123.4"))
	require.NoError(b, h.Storage.Update(context.Background(), "Dealloc", models.Gauge, "123.4"))

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(b, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAll)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
		assert.Equal(b, 200, rr.Code)
	}
}
