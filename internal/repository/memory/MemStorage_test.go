package memory

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"
	"os"
	"reflect"
	"testing"
)

func TestMemStorage_Update(t *testing.T) {
	memStorage := NewMemStorage("")

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
				valueType: models.Gauge,
			},
			err: nil,
		},
		{
			name:    "Counter int",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "2",
				valueType: models.Counter,
			},
			err: nil,
		},
		{
			name:    "Value is not number",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "hello",
				valueType: models.Counter,
			},
			err: apierror.NumberParse,
		},
		{
			name:    "Type and value mismatch",
			storage: memStorage,
			args: args{
				key:       "Alloc",
				value:     "123.245",
				valueType: models.Counter,
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
	storage := NewMemStorage("")
	require.NoError(t, storage.Update("Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update("PollCount", models.Counter, "2"))

	type args struct {
		key       string
		valueType models.MetricType
	}
	tests := []struct {
		name string
		args args
		want *models.Metric
		err  error
	}{
		{
			name: "Get metric",
			args: args{
				key:       "Alloc",
				valueType: models.Gauge,
			},
			want: utils.Ptr(models.Metric{
				Name:  "Alloc",
				MType: models.Gauge,
				Value: utils.Ptr(101.42),
			}),
			err: nil,
		},
		{
			name: "Get metric with another type",
			args: args{
				key:       "Alloc",
				valueType: models.Counter,
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
			got, err := storage.Get(tt.args.key, tt.args.valueType)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	storage := NewMemStorage("")

	require.NoError(t, storage.Update("Alloc", models.Gauge, "101.42"))
	require.NoError(t, storage.Update("PollCount", models.Counter, "2"))
	require.NoError(t, os.Setenv("KEY", ""))

	tests := []struct {
		name string
		want []models.Metric
	}{
		{
			name: "Get all",
			want: []models.Metric{
				{Name: "Alloc", MType: models.Gauge, Value: utils.Ptr(101.42)},
				{Name: "PollCount", MType: models.Counter, Delta: utils.Ptr(int64(2))},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, reflect.DeepEqual(tt.want, storage.GetAll()))
		})
	}
}
