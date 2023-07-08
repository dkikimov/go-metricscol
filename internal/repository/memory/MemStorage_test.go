package memory

import (
	"github.com/stretchr/testify/assert"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"testing"
)

func TestMemStorage_Update(t *testing.T) {
	memStorage := NewMemStorage()

	type args struct {
		key       string
		value     string
		valueType models.MetricType
	}
	tests := []struct {
		name    string
		storage *MemStorage
		args    args
		err     error
	}{
		{
			name:    "Gauge float",
			storage: memStorage,
			args: args{
				key:       "Alloc",
				value:     "120.123",
				valueType: models.GaugeType,
			},
			err: nil,
		},
		{
			name:    "Counter int",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "2",
				valueType: models.CounterType,
			},
			err: nil,
		},
		{
			name:    "Value is not number",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "hello",
				valueType: models.CounterType,
			},
			err: apierror.NumberParse,
		},
		{
			name:    "Type and value mismatch",
			storage: memStorage,
			args: args{
				key:       "Alloc",
				value:     "123.245",
				valueType: models.CounterType,
			},
			err: apierror.NumberParse,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assert.Equal(t, tt.err, tt.storage.Update(tt.args.key, tt.args.valueType, tt.args.value))
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	metrics := models.Metrics{}
	assert.NoError(t, metrics.Update("Alloc", models.GaugeType, 101.42))
	assert.NoError(t, metrics.Update("PollCount", models.CounterType, 2))

	type args struct {
		key       string
		valueType models.MetricType
	}
	tests := []struct {
		name string
		args args
		want models.Metric
		err  error
	}{
		{
			name: "Get metric",
			args: args{
				key:       "Alloc",
				valueType: models.GaugeType,
			},
			want: models.Gauge{
				Name:  "Alloc",
				Value: 101.42,
			},
			err: nil,
		},
		{
			name: "Get metric with another type",
			args: args{
				key:       "Alloc",
				valueType: models.CounterType,
			},
			want: nil,
			err:  apierror.NotFound,
		},
		{
			name: "Get metric with unknown type",
			args: args{
				key:       "Alloc",
				valueType: 5,
			},
			want: nil,
			err:  apierror.NotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := &MemStorage{
				metrics: metrics,
			}
			got, got1 := memStorage.Get(tt.args.key, tt.args.valueType)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, got1)
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	metrics := models.Metrics{}
	assert.NoError(t, metrics.Update("Alloc", models.GaugeType, 101.42))
	assert.NoError(t, metrics.Update("PollCount", models.CounterType, 2))

	type fields struct {
		metrics models.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		want   []models.Metric
	}{
		{
			name: "Get all",
			fields: fields{
				metrics: metrics,
			},
			want: []models.Metric{
				models.Gauge{Name: "Alloc", Value: 101.42},
				models.Counter{Name: "PollCount", Value: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := &MemStorage{
				metrics: tt.fields.metrics,
			}
			assert.EqualValues(t, tt.want, memStorage.GetAll())
		})
	}
}
