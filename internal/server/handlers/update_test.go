package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	"net/http"
	"net/http/httptest"
	"testing"
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
				metricType: models.CounterType.String(),
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Without value counter",
			args: args{
				metricType: models.CounterType.String(),
				metricName: "Alloc",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Without name and value gauge",
			args: args{
				metricType: models.GaugeType.String(),
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Without value gauge",
			args: args{
				metricType: models.GaugeType.String(),
				metricName: "Alloc",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Wrong value type",
			args: args{
				metricType:  models.GaugeType.String(),
				metricName:  "PollCount",
				metricValue: "hello",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Metric name and type mismatch counter",
			args: args{
				metricType:  models.CounterType.String(),
				metricName:  "Alloc",
				metricValue: "1.23",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handlers{
				Storage: memory.NewMemStorage(),
			}

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
