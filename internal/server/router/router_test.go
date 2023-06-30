package router

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

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
			args: args{addr: "/update/counter/PollCount/1.23"},
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

			statusCode, _ := testRequest(t, ts, http.MethodGet, tt.args.addr)
			assert.Equal(t, tt.want, statusCode)
		})
	}
}
