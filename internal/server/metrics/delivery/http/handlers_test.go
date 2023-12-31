package http

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
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/server"
	"go-metricscol/internal/server/metrics/usecase"
	"go-metricscol/internal/utils"
)

var emptyConfig, _ = server.NewConfig("", models.Duration{Duration: time.Second}, "", false, "", "", "", "")

func TestMetricsHandlers_Find(t *testing.T) {
	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		emptyConfig,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		emptyConfig,
	)

	require.NoError(t, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(123.4),
	}))

	require.NoError(t, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "MemoryInUse",
		MType: models.Gauge,
		Value: utils.Ptr(float64(593)),
	}))

	require.NoError(t, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "PollCount",
		MType: models.Counter,
		Delta: utils.Ptr(int64(1)),
	}))

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

func ExampleMetricsHandlers_Find() {
	address := "localhost:8080"

	metricType := models.Gauge
	metricName := "Alloc"

	findURL := fmt.Sprintf("%s/value/%s/%s", address, metricType, metricName)

	response, err := http.Get(findURL)
	if err != nil {
		// Handle error
	}
	response.Body.Close()
}

func TestMetricsHandlers_GetAll(t *testing.T) {
	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		nil,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		nil,
	)

	require.NoError(t, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(123.4),
	}))

	require.NoError(t, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "MemoryInUse",
		MType: models.Gauge,
		Value: utils.Ptr(float64(593)),
	}))

	require.NoError(t, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "PollCount",
		MType: models.Counter,
		Delta: utils.Ptr(int64(1)),
	}))

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

func ExampleMetricsHandlers_GetAll() {
	address := "localhost:8080"

	getAllURL := fmt.Sprintf("%s/", address)

	response, err := http.Get(getAllURL)
	if err != nil {
		// Handle error
	}
	response.Body.Close()
}

func TestMetricsHandlers_GetAllWithHash(t *testing.T) {
	hashKey := "test"
	config, _ := server.NewConfig("", models.Duration{Duration: time.Second}, "", false, hashKey, "", "", "")
	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		config,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		config,
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

	require.NoError(t, h.metricsUC.Update(context.Background(), alloc))
	require.NoError(t, h.metricsUC.Update(context.Background(), memoryInUse))
	require.NoError(t, h.metricsUC.Update(context.Background(), pollCount))

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

func TestMetricsHandlers_FindJSON(t *testing.T) {
	storage := memory.NewMemStorage()
	require.NoError(t, storage.Update(context.Background(), models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(12.1),
	}))
	require.NoError(t, storage.Update(context.Background(), models.Metric{
		Name:  "PollCount",
		MType: models.Counter,
		Delta: utils.Ptr(int64(13)),
	}))

	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		emptyConfig,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		emptyConfig,
	)

	type want struct {
		Body       models.Metric
		StatusCode int
	}

	tests := []struct {
		name string
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

func ExampleMetricsHandlers_FindJSON() {
	address := "localhost:8080"

	metricToFind := models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
	}

	marshaledMetric, err := json.Marshal(metricToFind)
	if err != nil {
		// Handle error
		return
	}

	updatePostURL := fmt.Sprintf("%s/value/", address)

	response, err := http.Post(updatePostURL, "application/json", bytes.NewReader(marshaledMetric))
	if err != nil {
		// Handle error
		return
	}
	response.Body.Close()
}

func BenchmarkMetricsHandlers_Find_MemStorage(b *testing.B) {
	config, _ := server.NewConfig("", models.Duration{Duration: time.Second}, "", false, "hash", "", "", "")

	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		config,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		config,
	)

	require.NoError(b, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(123.4),
	}))

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

func BenchmarkMetricsHandlers_FindJSON_MemStorage(b *testing.B) {
	config, _ := server.NewConfig("", models.Duration{Duration: time.Second}, "", false, "hash", "", "", "")

	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		config,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		config,
	)

	metric := models.Metric{Name: "Alloc", MType: models.Gauge, Value: utils.Ptr(123.4)}
	require.NoError(b, h.metricsUC.Update(context.Background(), metric))

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

func BenchmarkMetricsHandlers_FindAllWithHash_MemStorage(b *testing.B) {
	config, _ := server.NewConfig("", models.Duration{Duration: time.Second}, "", false, "hash", "", "", "")

	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		config,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		config,
	)

	require.NoError(b, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: utils.Ptr(123.4),
	}))

	require.NoError(b, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "MemoryInUse",
		MType: models.Gauge,
		Value: utils.Ptr(float64(593)),
	}))

	require.NoError(b, h.metricsUC.Update(context.Background(), models.Metric{
		Name:  "PollCount",
		MType: models.Counter,
		Delta: utils.Ptr(int64(1)),
	}))

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

func TestMetricsHandlers_Update(t *testing.T) {
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
			newMetricsUC := usecase.NewMetricsUC(
				memory.NewMemStorage(),
				nil,
				emptyConfig,
			)
			h := newMetricsHandlers(
				newMetricsUC,
				emptyConfig,
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

func TestMetricsHandlers_UpdateJSON(t *testing.T) {
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

		newMetricsUC := usecase.NewMetricsUC(
			storage,
			nil,
			emptyConfig,
		)
		h := newMetricsHandlers(
			newMetricsUC,
			emptyConfig,
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

func ExampleMetricsHandlers_Update() {
	address := "localhost:8080"
	metricType := models.Gauge
	metricName := "Alloc"
	metricValue := 5

	updatePostURL := fmt.Sprintf("%s/update/%s/%s/%d", address, metricType, metricName, metricValue)

	response, err := http.Post(updatePostURL, "text/plain", nil)
	if err != nil {
		// Handle error
		return
	}
	response.Body.Close()
}

func BenchmarkMetricsHandlers_Update_MemStorage(b *testing.B) {
	config, _ := server.NewConfig("", models.Duration{Duration: time.Second}, "", false, "hash", "", "", "")

	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		config,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		config,
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

func ExampleMetricsHandlers_UpdateJSON() {
	address := "localhost:8080"

	var metricValue float64 = 1
	metricToUpdate := models.Metric{
		Name:  "Alloc",
		MType: models.Gauge,
		Value: &metricValue,
	}
	marshaledMetric, err := json.Marshal(metricToUpdate)
	if err != nil {
		// Handle error
		return
	}

	updatePostURL := fmt.Sprintf("%s/update/", address)

	response, err := http.Post(updatePostURL, "application/json", bytes.NewReader(marshaledMetric))
	if err != nil {
		// Handle error
		return
	}
	response.Body.Close()
}

func BenchmarkMetricsHandlers_UpdateJSON_MemStorage(b *testing.B) {
	config, _ := server.NewConfig("", models.Duration{Duration: time.Second}, "", false, "hash", "", "", "")

	newMetricsUC := usecase.NewMetricsUC(
		memory.NewMemStorage(),
		nil,
		config,
	)
	h := newMetricsHandlers(
		newMetricsUC,
		config,
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

func TestMetricsHandlers_Updates(t *testing.T) {
	address := "localhost:8080"

	var allocValue float64 = 1
	var countValue int64 = 2

	metricsToUpdate := []models.Metric{
		{
			Name:  "Alloc",
			MType: models.Gauge,
			Value: &allocValue,
		},
		{
			Name:  "Count",
			MType: models.Counter,
			Delta: &countValue,
		},
	}
	marshaledMetrics, _ := json.Marshal(metricsToUpdate)

	updatePostURL := fmt.Sprintf("%s/updates/", address)

	response, err := http.Post(updatePostURL, "application/json", bytes.NewReader(marshaledMetrics))
	if err != nil {
		// Handle error
		return
	}
	response.Body.Close()
}
