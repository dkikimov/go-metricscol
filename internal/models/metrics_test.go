package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/server/apierror"
	"testing"
)

func TestMetrics_Get(t *testing.T) {
	metrics := Metrics{}
	metrics.Update("Alloc", GaugeType, 13)
	metrics.Update("PollCount", CounterType, 1)

	type args struct {
		name      string
		valueType MetricType
	}
	tests := []struct {
		name string
		args args
		want Metric
		err  error
	}{
		{
			name: "Get metric gauge",
			args: args{
				name:      "Alloc",
				valueType: GaugeType,
			},
			want: Gauge{
				Name:  "Alloc",
				Value: 13,
			},
			err: nil,
		},
		{
			name: "Get metric counter",
			args: args{
				name:      "PollCount",
				valueType: CounterType,
			},
			want: Counter{
				Name:  "PollCount",
				Value: 1,
			},
			err: nil,
		},
		{
			name: "Metric not found",
			args: args{
				name:      "La",
				valueType: GaugeType,
			},
			want: nil,
			err:  apierror.NotFound,
		},
		{
			name: "Metric with another type",
			args: args{
				name:      "Alloc",
				valueType: CounterType,
			},
			want: nil,
			err:  apierror.NotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := metrics.Get(tt.args.name, tt.args.valueType)

			assert.Equal(t, tt.want, got)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestMetrics_ResetPollCount(t *testing.T) {
	t.Run("Reset poll count", func(t *testing.T) {
		metrics := Metrics{}

		require.NoError(t, metrics.Update("PollCount", CounterType, 2))

		metrics.ResetPollCount()
		metric, err := metrics.Get("PollCount", CounterType)

		assert.EqualValues(t, nil, err)
		assert.EqualValues(t, 0, metric.(Counter).Value)
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
		want Metric
		err  error
	}{
		{
			name: "Update gauge int",
			args: args{
				name:      "Alloc",
				valueType: GaugeType,
				value:     13,
			},
			want: Gauge{
				Name:  "Alloc",
				Value: 13,
			},
			err: nil,
		},
		{
			name: "Update gauge float",
			args: args{
				name:      "Alloc",
				valueType: GaugeType,
				value:     13.41,
			},
			want: Gauge{
				Name:  "Alloc",
				Value: 13.41,
			},
			err: nil,
		},
		{
			name: "Update counter int",
			args: args{
				name:      "PollCount",
				valueType: CounterType,
				value:     1,
			},
			want: Counter{
				Name:  "PollCount",
				Value: 1,
			},
			err: nil,
		},
		{
			name: "Unknown type",
			args: args{
				name:      "PollCount",
				valueType: 4,
				value:     1,
			},
			want: nil,
			err:  apierror.UnknownMetricType,
		},

		{
			name: "Invalid value counter",
			args: args{
				name:      "PollCount",
				valueType: CounterType,
				value:     1.34,
			},
			want: nil,
			err:  apierror.InvalidValue,
		},
		{
			name: "Invalid value gauge",
			args: args{
				name:      "Alloc",
				valueType: GaugeType,
				value:     "",
			},
			want: nil,
			err:  apierror.InvalidValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metrics{}
			err := m.Update(tt.args.name, tt.args.valueType, tt.args.value)

			require.EqualValues(t, tt.err, err)
			require.Equal(t, tt.want, m[getKey(tt.args.name, tt.args.valueType)])
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
				valueType: GaugeType,
			},
			want: "Alloc:gauge",
		},
		{
			name: "Get key counter",
			args: args{
				name:      "PollCount",
				valueType: CounterType,
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
