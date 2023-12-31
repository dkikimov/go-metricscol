package middleware

import (
	"go-metricscol/internal/server"
	"go-metricscol/internal/server/metrics"
)

type Manager struct {
	metricsUC metrics.UseCase
	cfg       *server.Config
	origins   []string
}

func NewMiddlewareManager(metricsUC metrics.UseCase, cfg *server.Config, origins []string) *Manager {
	return &Manager{metricsUC: metricsUC, cfg: cfg, origins: origins}
}
