package repository

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
)

type Repository interface {
	Update(key string, value string, valueType models.MetricType) apierror.APIError
}
