package handlers

import (
	"go-metricscol/internal/repository/memory"
	"net/http"
	"net/http/httptest"
	"testing"
)

const baseURL = "http://127.0.0.1:8080"

func TestHandlers(t *testing.T) {
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
			args: args{addr: baseURL + "/update/counter/"},
			want: http.StatusNotFound,
		},
		{
			name: "Without value counter",
			args: args{addr: baseURL + "/update/counter/Alloc"},
			want: http.StatusNotFound,
		},
		{
			name: "Without name and value gauge",
			args: args{addr: baseURL + "/update/gauge/"},
			want: http.StatusNotFound,
		},
		{
			name: "Without value gauge",
			args: args{addr: baseURL + "/update/gauge/Alloc"},
			want: http.StatusNotFound,
		},
		{
			name: "Wrong value type",
			args: args{addr: baseURL + "/update/counter/PollCount/1.23"},
			want: http.StatusBadRequest,
		},
		{
			name: "Metric name and type mismatch counter",
			args: args{addr: baseURL + "/update/counter/Alloc/1.23"},
			want: http.StatusBadRequest,
		},
	}

	processors := Processors{Storage: memory.NewMemStorage()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.args.addr, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(processors.Update)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want {
				t.Errorf("Expected status code %d, got %d", tt.want, w.Code)
			}
			res.Body.Close()
		})
	}
}
