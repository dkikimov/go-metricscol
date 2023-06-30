package memory

import (
	"github.com/stretchr/testify/assert"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"net/http"
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
				valueType: models.Gauge,
			},
			want: http.StatusOK,
		},
		{
			name:    "Counter int",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "2",
				valueType: models.Counter,
			},
			want: http.StatusOK,
		},
		{
			name:    "Value is not number",
			storage: memStorage,
			args: args{
				key:       "PollCount",
				value:     "hello",
				valueType: models.Counter,
			},
			want: http.StatusBadRequest,
		},
		{
			name:    "Type and value mismatch",
			storage: memStorage,
			args: args{
				key:       "Alloc",
				value:     "123.245",
				valueType: models.Counter,
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.storage.Update(tt.args.key, tt.args.value, tt.args.valueType), tt.want)
		})
	}
}
