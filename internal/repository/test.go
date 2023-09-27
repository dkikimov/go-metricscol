package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/utils"
)

func TestUpdate(ctx context.Context, t *testing.T, storage Repository) {
	type args struct {
		key       string
		value     string
		valueType models.MetricType
	}
	tests := []struct {
		name    string
		storage Repository
		args    args
		err     error
	}{
		{
			name:    "Gauge float",
			storage: storage,
			args: args{
				key:       "Alloc",
				value:     "120.123",
				valueType: models.Gauge,
			},
			err: nil,
		},
		{
			name:    "Counter int",
			storage: storage,
			args: args{
				key:       "PollCount",
				value:     "2",
				valueType: models.Counter,
			},
			err: nil,
		},
		{
			name:    "Value is not number",
			storage: storage,
			args: args{
				key:       "PollCount",
				value:     "hello",
				valueType: models.Counter,
			},
			err: apierror.NumberParse,
		},
		{
			name:    "Type and value mismatch",
			storage: storage,
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
			assert.Equal(t, tt.err, tt.storage.Update(ctx, tt.args.key, tt.args.valueType, tt.args.value))
		})
	}
}

func TestGet(ctx context.Context, t *testing.T, storage Repository) {
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
			name: "Find metric gauge",
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
			name: "Find metric counter",
			args: args{
				key:       "PollCount",
				valueType: models.Counter,
			},
			want: utils.Ptr(models.Metric{
				Name:  "PollCount",
				MType: models.Counter,
				Delta: utils.Ptr(int64(1)),
			}),
			err: nil,
		},
		{
			name: "Find metric with another type",
			args: args{
				key:       "Alloc",
				valueType: models.Counter,
			},
			want: nil,
			err:  apierror.NotFound,
		},
		{
			name: "Find metric with unknown type",
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
			got, err := storage.Get(ctx, tt.args.key, tt.args.valueType)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGetAll(ctx context.Context, t *testing.T, storage Repository) {
	tests := []struct {
		name string
		want []models.Metric
	}{
		{
			name: "Find all",
			want: []models.Metric{
				{Name: "Alloc", MType: models.Gauge, Value: utils.Ptr(101.42)},
				{Name: "PollCount", MType: models.Counter, Delta: utils.Ptr(int64(1))},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			all, err := storage.GetAll(ctx)
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tt.want, all))
		})
	}
}

func TestUpdateWithStruct(ctx context.Context, t *testing.T, storage Repository) {
	tests := []struct {
		name    string
		storage Repository
		args    models.Metric
		err     error
	}{
		{
			name:    "Gauge float",
			storage: storage,
			args: models.Metric{
				Name:  "Alloc",
				Value: utils.Ptr(120.123),
				MType: models.Gauge,
			},
			err: nil,
		},
		{
			name:    "Counter int",
			storage: storage,
			args: models.Metric{
				Name:  "PollCount",
				Delta: utils.Ptr(int64(2)),
				MType: models.Counter,
			},
			err: nil,
		},
		{
			name:    "Type and value mismatch",
			storage: storage,
			args: models.Metric{
				Name:  "Alloc",
				Value: utils.Ptr(123.245),
				MType: models.Counter,
			},
			err: apierror.InvalidValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.UpdateWithStruct(ctx, &tt.args)
			assert.Equal(t, tt.err, err)
			if err == nil {
				m, err := storage.Get(ctx, tt.args.Name, tt.args.MType)
				require.NoError(t, err)
				assert.Equal(t, tt.args, *m)
			}
		})
	}
}

func TestUpdates(ctx context.Context, t *testing.T, storage Repository) {
	tests := []struct {
		name    string
		storage Repository
		args    []models.Metric
		err     error
	}{
		{
			name:    "Updates",
			storage: storage,
			args: []models.Metric{
				{
					Name:  "Alloc",
					Value: utils.Ptr(120.123),
					MType: models.Gauge,
				},
				{
					Name:  "PollCount",
					Delta: utils.Ptr(int64(1)),
					MType: models.Counter,
				},
			},
			err: nil,
		},
		{
			name:    "Gauge with delta",
			storage: storage,
			args: []models.Metric{
				{
					Name:  "Alloc",
					Value: utils.Ptr(120.123),
					MType: models.Gauge,
				},
				{
					Name:  "PollCount",
					Delta: utils.Ptr(int64(1)),
					MType: models.Gauge,
				},
			},
			err: apierror.InvalidValue,
		},
		{
			name:    "Counter with non-empty value field",
			storage: storage,
			args: []models.Metric{
				{
					Name:  "Alloc",
					Value: utils.Ptr(120.123),
					MType: models.Gauge,
				},
				{
					Name:  "PollCount",
					Value: utils.Ptr(1.34),
					MType: models.Counter,
				},
			},
			err: apierror.InvalidValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.Updates(ctx, tt.args)
			assert.Equal(t, tt.err, err)

			all, err := storage.GetAll(ctx)
			require.NoError(t, err)

			if tt.err == nil {
				assert.EqualValues(t, tt.args, all)
			} else {
				if ok := storage.SupportsTx(); ok {
					assert.Empty(t, all)
				}
			}
		})
	}
}
