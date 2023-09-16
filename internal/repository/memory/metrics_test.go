package memory

import (
	"reflect"
	"testing"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_Get(t *testing.T) {
	metrics := NewMetrics()

	require.NoError(t, metrics.Update("Alloc", models.Gauge, 13))
	require.NoError(t, metrics.Update("PollCount", models.Counter, 1))

	type args struct {
		name      string
		valueType models.MetricType
	}
	tests := []struct {
		name string
		args args
		want *models.Metric
		err  error
	}{
		{
			name: "Find metric gauge",
			args: args{
				name:      "Alloc",
				valueType: models.Gauge,
			},
			want: utils.Ptr(models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
				Value: utils.Ptr(float64(13)),
			}),
			err: nil,
		},
		{
			name: "Find metric counter",
			args: args{
				name:      "PollCount",
				valueType: models.Counter,
			},
			want: utils.Ptr(models.Metric{
				Name:  "PollCount",
				MType: models.Counter,
				Delta: utils.Ptr(int64(13)),
			}),
			err: nil,
		},
		{
			name: "Metric not found",
			args: args{
				name:      "La",
				valueType: models.Gauge,
			},
			want: nil,
			err:  apierror.NotFound,
		},
		{
			name: "Metric with another type",
			args: args{
				name:      "Alloc",
				valueType: models.Counter,
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

		require.NoError(t, metrics.Update("PollCount", models.Counter, 2))

		metrics.ResetPollCount()
		metric, err := metrics.Get("PollCount", models.Counter)

		assert.EqualValues(t, nil, err)
		if err == nil {
			assert.True(t, reflect.DeepEqual(metric.Delta, utils.Ptr(int64(0))))
		}
	})
}

func TestMetrics_Update(t *testing.T) {
	type args struct {
		name      string
		valueType models.MetricType
		value     interface{}
	}

	tests := []struct {
		name string
		args args
		want *models.Metric
		err  error
	}{
		{
			name: "Update gauge int",
			args: args{
				name:      "Alloc",
				valueType: models.Gauge,
				value:     13,
			},
			want: utils.Ptr(models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
				Value: utils.Ptr(float64(13)),
			}),
			err: nil,
		},
		{
			name: "Update gauge float",
			args: args{
				name:      "Alloc",
				valueType: models.Gauge,
				value:     13.41,
			},
			want: utils.Ptr(models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
				Value: utils.Ptr(13.41),
			}),
			err: nil,
		},
		{
			name: "Update counter int",
			args: args{
				name:      "PollCount",
				valueType: models.Counter,
				value:     1,
			},
			want: utils.Ptr(models.Metric{
				Name:  "PollCount",
				MType: models.Counter,
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
				valueType: models.Counter,
				value:     1.34,
			},
			want: nil,
			err:  apierror.InvalidValue,
		},
		{
			name: "Invalid value gauge",
			args: args{
				name:      "Alloc",
				valueType: models.Gauge,
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
				assert.True(t, true, reflect.DeepEqual(m.Collection[getKey(tt.args.name, tt.args.valueType)], tt.want))
			}
		})
	}
}

func Test_getKey(t *testing.T) {
	type args struct {
		name      string
		valueType models.MetricType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Find key gauge",
			args: args{
				name:      "Alloc",
				valueType: models.Gauge,
			},
			want: "Allocg",
		},
		{
			name: "Find key counter",
			args: args{
				name:      "PollCount",
				valueType: models.Counter,
			},
			want: "PollCountc",
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
		metric  models.Metric
		wantErr error
	}{
		{
			name: "Update gauge int",
			metric: models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
				Value: utils.Ptr(float64(13)),
			},
			wantErr: nil,
		},
		{
			name: "Update gauge float",
			metric: models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
				Value: utils.Ptr(13.41),
			},
			wantErr: nil,
		},
		{
			name: "Update counter int",
			metric: models.Metric{
				Name:  "PollCount",
				MType: models.Counter,
				Delta: utils.Ptr(int64(1)),
			},
			wantErr: nil,
		},
		{
			name: "Unknown type",
			metric: models.Metric{
				Name:  "PollCount",
				MType: "unknown",
				Delta: utils.Ptr(int64(1)),
			},
			wantErr: apierror.UnknownMetricType,
		},
		{
			name: "Delta for gauge",
			metric: models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
				Delta: utils.Ptr(int64(1)),
			},
			wantErr: apierror.InvalidValue,
		},
		{
			name: "Value for counter",
			metric: models.Metric{
				Name:  "PollCount",
				MType: models.Counter,
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
