package router

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: Какие тесты нужно писать для роутера? Как тестить хендлеры в своей директории, избегая циклического импорта?

func TestHandlers_Get(t *testing.T) {
	storage := memory.NewMemStorage()
	storage.Update("Alloc", models.GaugeType, "123.4")
	storage.Update("MemoryInUse", models.GaugeType, "593")
	storage.Update("PollCount", models.CounterType, "1")

	type fields struct {
		Storage repository.Repository
	}
	type want struct {
		StatusCode int
		Body       string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "Get Alloc value",
			fields: fields{Storage: storage},
			args: args{
				url: "/value/gauge/Alloc",
			},
			want: want{
				StatusCode: http.StatusOK,
				Body:       "123.4",
			},
		},
		{
			name:   "Unknown metric",
			fields: fields{Storage: storage},
			args: args{
				url: "/value/gauge/NewMetric",
			},
			want: want{
				StatusCode: http.StatusNotFound,
				Body:       "",
			},
		},
		{
			name:   "Unknown metric type",
			fields: fields{Storage: storage},
			args: args{
				url: "/value/h/Alloc",
			},
			want: want{
				StatusCode: http.StatusNotImplemented,
				Body:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewWithStorage(storage)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, body := utils.TestRequest(t, ts, http.MethodGet, tt.args.url)

			require.Equal(t, tt.want.StatusCode, statusCode)
			require.Equal(t, tt.want.Body, body)
		})
	}
}

func TestHandlers_GetAll(t *testing.T) {
	storage := memory.NewMemStorage()
	storage.Update("Alloc", models.GaugeType, "123.4")
	storage.Update("MemoryInUse", models.GaugeType, "593")
	storage.Update("PollCount", models.CounterType, "1")

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
			r := NewWithStorage(storage)
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, body := utils.TestRequest(t, ts, http.MethodGet, tt.args.url)

			require.Equal(t, tt.want.StatusCode, statusCode)
			require.Equal(t, tt.want.Body, body)
		})
	}
}

func TestHandlers_Update(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want int // status code
	}{
		{
			name: "Without name and value counter",
			args: args{addr: "/update/counter/"},
			want: http.StatusNotFound,
		},
		{
			name: "Without value counter",
			args: args{addr: "/update/counter/Alloc"},
			want: http.StatusNotFound,
		},
		{
			name: "Without name and value gauge",
			args: args{addr: "/update/gauge/"},
			want: http.StatusNotFound,
		},
		{
			name: "Without value gauge",
			args: args{addr: "/update/gauge/Alloc"},
			want: http.StatusNotFound,
		},
		{
			name: "Wrong value type",
			args: args{addr: "/update/counter/PollCount/hello"},
			want: http.StatusBadRequest,
		},
		{
			name: "Metric name and type mismatch counter",
			args: args{addr: "/update/counter/Alloc/1.23"},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New()
			ts := httptest.NewServer(r)
			defer ts.Close()

			statusCode, _ := utils.TestRequest(t, ts, http.MethodPost, tt.args.addr)
			assert.Equal(t, tt.want, statusCode)
		})
	}
}
