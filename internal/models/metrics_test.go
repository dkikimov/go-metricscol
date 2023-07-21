package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"
	"os"
	"reflect"
	"testing"
)

func TestMetrics_Get(t *testing.T) {
	metrics := NewMetrics()

	require.NoError(t, metrics.Update("Alloc", Gauge, 13))
	require.NoError(t, metrics.Update("PollCount", Counter, 1))

	type args struct {
		name      string
		valueType MetricType
	}
	tests := []struct {
		name string
		args args
		want *Metric
		err  error
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
			err: nil,
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
			err: nil,
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
			if err == nil {
				assert.True(t, reflect.DeepEqual(tt.want.Value, got.Value))
			}
		})
	}
}

func TestMetrics_ResetPollCount(t *testing.T) {
	t.Run("Reset poll count", func(t *testing.T) {
		metrics := NewMetrics()

		require.NoError(t, metrics.Update("PollCount", Counter, 2))

		metrics.ResetPollCount()
		metric, err := metrics.Get("PollCount", Counter)

		assert.EqualValues(t, nil, err)
		if err == nil {
			assert.True(t, reflect.DeepEqual(metric.Delta, utils.Ptr(int64(0))))
		}
	})
}

func TestMetrics_Update(t *testing.T) {
	type args struct {
		name      string
		valueType MetricType
		value     interface{}
	}

	key := "test"
	require.NoError(t, os.Setenv("KEY", key))

	tests := []struct {
		name string
		args args
		want *Metric
		err  error
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
			err: nil,
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
			err: nil,
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
			err: nil,
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
			m := NewMetrics()
			err := m.Update(tt.args.name, tt.args.valueType, tt.args.value)

			assert.EqualValues(t, tt.err, err)
			if err == nil {
				tt.want.SetHashValue(key)

				assert.True(t, true, reflect.DeepEqual(m.Collection[getKey(tt.args.name, tt.args.valueType)], tt.want))
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

func TestMetrics_UpdateWithStruct(t *testing.T) {
	tests := []struct {
		name    string
		metric  Metric
		wantErr error
	}{
		{
			name: "Update gauge int",
			metric: Metric{
				Name:  "Alloc",
				MType: Gauge,
				Value: utils.Ptr(float64(13)),
			},
			wantErr: nil,
		},
		{
			name: "Update gauge float",
			metric: Metric{
				Name:  "Alloc",
				MType: Gauge,
				Value: utils.Ptr(13.41),
			},
			wantErr: nil,
		},
		{
			name: "Update counter int",
			metric: Metric{
				Name:  "PollCount",
				MType: Counter,
				Delta: utils.Ptr(int64(1)),
			},
			wantErr: nil,
		},
		{
			name: "Unknown type",
			metric: Metric{
				Name:  "PollCount",
				MType: "unknown",
				Delta: utils.Ptr(int64(1)),
			},
			wantErr: apierror.UnknownMetricType,
		},
		{
			name: "Delta for gauge",
			metric: Metric{
				Name:  "Alloc",
				MType: Gauge,
				Delta: utils.Ptr(int64(1)),
			},
			wantErr: apierror.InvalidValue,
		},
		{
			name: "Value for counter",
			metric: Metric{
				Name:  "PollCount",
				MType: Counter,
				Value: utils.Ptr(float64(1)),
			},
			wantErr: apierror.InvalidValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()

			err := m.UpdateWithStruct(&tt.metric)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				assert.True(t, true, reflect.DeepEqual(m.Collection[getKey(tt.metric.Name, tt.metric.MType)], tt.metric))
			}
		})
	}
}
