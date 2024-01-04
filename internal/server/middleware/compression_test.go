package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-metricscol/internal/config"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
	helathUseCase "go-metricscol/internal/server/health/usecase"
	metricsUseCase "go-metricscol/internal/server/metrics/usecase"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Test_decompressHandler(t *testing.T) {
	type args struct {
		text   string
		isGzip bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Usual request",
			args: args{
				text:   "Hello world!",
				isGzip: false,
			},
		},
		{
			name: "Gzip request",
			args: args{
				text:   "Hello gzipped world!",
				isGzip: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBuffer([]byte{})

			if tt.args.isGzip {
				writer := gzip.NewWriter(body)
				_, err := writer.Write([]byte(tt.args.text))
				require.NoError(t, err)
				require.NoError(t, writer.Close())
			} else {
				body.WriteString(tt.args.text)
			}

			req, err := http.NewRequest(http.MethodPost, "", body)
			require.NoError(t, err)

			if tt.args.isGzip {
				req.Header.Set("Content-Encoding", "gzip")
			}

			rr := httptest.NewRecorder()

			repository := memory.NewMemStorage()
			cfg, err := config.NewServerConfig("", models.Duration{Duration: time.Second}, "", false, "", "", "", "")
			require.NoError(t, err)

			metricsUC := metricsUseCase.NewMetricsUC(repository, nil, cfg)
			healthUC := helathUseCase.NewHealthUC(nil)

			mw := NewManager(metricsUC, healthUC, cfg, repository)

			handler := mw.DecompressHandler(http.HandlerFunc(testHandler))

			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, tt.args.text, rr.Body.String())
		})

	}
}
