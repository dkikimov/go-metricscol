package memory

import (
	"github.com/stretchr/testify/assert"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"
	"net/http"
	"reflect"
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
		want    apierror.APIError
	}{
		{
			name:    "Gauge float",
			storage: memStorage,
			args: args{
				key:       "Alloc",
				value:     "120.123",
				valueType: models.GaugeType,
			},
			want: http.StatusOK,
		},
		{
			name:    "Counter int",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "2",
				valueType: models.CounterType,
			},
			want: http.StatusOK,
		},
		{
			name:    "Value is not number",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "hello",
				valueType: models.CounterType,
			},
			want: http.StatusBadRequest,
		},
		{
			name:    "Type and value mismatch",
			storage: memStorage,
			args: args{
				key:       "Alloc",
				value:     "123.245",
				valueType: models.CounterType,
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, reflect.DeepEqual(tt.storage.Update(tt.args.key, tt.args.valueType, tt.args.value), tt.want))
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	metrics := models.MetricsMap{}
	metrics.Update("Alloc", models.GaugeType, 101.42)
	metrics.Update("PollCount", models.CounterType, 2)

	type args struct {
		key       string
		valueType models.MetricType
	}
	tests := []struct {
		name string
		args args
		want *models.Metric
		err  apierror.APIError
	}{
		{
			name: "Get metric",
			args: args{
				key:       "Alloc",
				valueType: models.GaugeType,
			},
			want: utils.Ptr(models.Metric{
				Name:  "Alloc",
				MType: models.GaugeType,
				Value: utils.Ptr(101.42),
			}),
			err: apierror.NoError,
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
				valueType: "unknown",
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
	metrics := models.MetricsMap{}
	metrics.Update("Alloc", models.GaugeType, 101.42)
	metrics.Update("PollCount", models.CounterType, 2)

	type fields struct {
		metrics models.MetricsMap
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
				{Name: "Alloc", MType: models.GaugeType, Value: utils.Ptr(101.42)},
				{Name: "PollCount", MType: models.CounterType, Delta: utils.Ptr(int64(2))},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := &MemStorage{
				metrics: tt.fields.metrics,
			}

			assert.True(t, reflect.DeepEqual(tt.want, memStorage.GetAll()))
		})
	}
}
