package repository

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
)

type Repository interface {
	Update(name string, valueType models.MetricType, value string) apierror.APIError
	Get(name string, valueType models.MetricType) (*models.Metric, apierror.APIError)
	GetAll() []models.Metric
}
