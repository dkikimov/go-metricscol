package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"
	"reflect"
	"testing"
)

func TestMetrics_Get(t *testing.T) {
	metrics := MetricsMap{}
	metrics.Update("Alloc", Gauge, 13)
	metrics.Update("PollCount", Counter, 1)

	type args struct {
		name      string
		valueType MetricType
	}
	tests := []struct {
		name string
		args args
		want *Metric
		err  apierror.APIError
	}{
		{
			name: "Get metric gauge",
			args: args{
				name:      "Alloc",
				valueType: Gauge,
			},
			want: utils.Ptr(Metric{
				Name:  "Alloc",
				MType: Gauge,
				Value: utils.Ptr(float64(13)),
			}),
			err: apierror.NoError,
		},
		{
			name: "Get metric counter",
			args: args{
				name:      "PollCount",
				valueType: Counter,
			},
			want: utils.Ptr(Metric{
				Name:  "PollCount",
				MType: Counter,
				Delta: utils.Ptr(int64(13)),
			}),
			err: apierror.NoError,
		},
		{
			name: "Metric not found",
			args: args{
				name:      "La",
				valueType: Gauge,
			},
			want: nil,
			err:  apierror.NotFound,
		},
		{
			name: "Metric with another type",
			args: args{
				name:      "Alloc",
				valueType: Counter,
			},
			want: nil,
			err:  apierror.NotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := metrics.Get(tt.args.name, tt.args.valueType)

			assert.EqualValues(t, tt.err, err)
			if tt.err == apierror.NoError {
				assert.True(t, true, reflect.DeepEqual(tt.want.Value, got.Value))
			}
		})
	}
}

func TestMetrics_ResetPollCount(t *testing.T) {
	metrics := MetricsMap{}
	metrics.Update("PollCount", Counter, 2)
	t.Run("Reset poll count", func(t *testing.T) {
		metrics.ResetPollCount()
		metric, err := metrics.Get("PollCount", Counter)

		assert.EqualValues(t, apierror.NoError, err)
		assert.True(t, reflect.DeepEqual(metric.Delta, utils.Ptr(int64(0))))
	})
}

func TestMetrics_Update(t *testing.T) {
	type args struct {
		name      string
		valueType MetricType
		value     interface{}
	}
	tests := []struct {
		name string
		args args
		want *Metric
		err  apierror.APIError
	}{
		{
			name: "Update gauge int",
			args: args{
				name:      "Alloc",
				valueType: Gauge,
				value:     13,
			},
			want: utils.Ptr(Metric{
				Name:  "Alloc",
				MType: Gauge,
				Value: utils.Ptr(float64(13)),
			}),
			err: apierror.NoError,
		},
		{
			name: "Update gauge float",
			args: args{
				name:      "Alloc",
				valueType: Gauge,
				value:     13.41,
			},
			want: utils.Ptr(Metric{
				Name:  "Alloc",
				MType: Gauge,
				Value: utils.Ptr(13.41),
			}),
			err: apierror.NoError,
		},
		{
			name: "Update counter int",
			args: args{
				name:      "PollCount",
				valueType: Counter,
				value:     1,
			},
			want: utils.Ptr(Metric{
				Name:  "PollCount",
				MType: Counter,
				Delta: utils.Ptr(int64(1)),
			}),
			err: apierror.NoError,
		},
		{
			name: "Unknown type",
			args: args{
				name:      "PollCount",
				valueType: "unknown",
				value:     1,
			},
			want: nil,
			err:  apierror.UnknownMetricType,
		},

		{
			name: "Invalid value counter",
			args: args{
				name:      "PollCount",
				valueType: Counter,
				value:     1.34,
			},
			want: nil,
			err:  apierror.InvalidValue,
		},
		{
			name: "Invalid value gauge",
			args: args{
				name:      "Alloc",
				valueType: Gauge,
				value:     "",
			},
			want: nil,
			err:  apierror.InvalidValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MetricsMap{}
			err := m.Update(tt.args.name, tt.args.valueType, tt.args.value)

			require.EqualValues(t, tt.err, err)
			if err == apierror.NoError {
				require.True(t, true, reflect.DeepEqual(m[getKey(tt.args.name, tt.args.valueType)], tt.want))

			}
		})
	}
}

func Test_getKey(t *testing.T) {
	type args struct {
		name      string
		valueType MetricType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get key gauge",
			args: args{
				name:      "Alloc",
				valueType: Gauge,
			},
			want: "Alloc:gauge",
		},
		{
			name: "Get key counter",
			args: args{
				name:      "PollCount",
				valueType: Counter,
			},
			want: "PollCount:counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getKey(tt.args.name, tt.args.valueType))
		})
	}
}
