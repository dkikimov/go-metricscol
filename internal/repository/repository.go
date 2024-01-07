package repository

import (
	"context"

	"go-metricscol/internal/models"
)

// Repository is interface that describes the storage of models.Metric.
type Repository interface {
	// Update adds or replaces existing metric with new one.
	Update(ctx context.Context, metric models.Metric) error

	// Updates adds or replaces multiple metrics in storage.
	Updates(ctx context.Context, metrics []models.Metric) error

	// UpdateWithStruct adds or replaces metric that was passed as models.Metric struct.
	UpdateWithStruct(ctx context.Context, metric *models.Metric) error

	// Get returns models.Metric if found.
	// If not apierror.NotFound error and nil models.Metric pointer returned.
	Get(ctx context.Context, key string, valueType models.MetricType) (*models.Metric, error)

	// GetAll returns slice of all models.Metric stored in repository.
	GetAll(ctx context.Context) ([]models.Metric, error)

	// SupportsTx returns if repository supports transactions.
	SupportsTx() bool

	// SupportsSavingToDisk returns if repository supports saving to disk.
	SupportsSavingToDisk() bool

	// SaveToDisk saves contents of repository to disk.
	SaveToDisk(filePath string) error

	RestoreFromDisk(filePath string) error

	// Ping returns no error if connection to repository is alive.
	Ping(ctx context.Context) error
}
