package middleware

import (
	"go-metricscol/internal/config"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/health"
	"go-metricscol/internal/server/metrics"
)

type Manager struct {
	metricsUC metrics.UseCase
	healthUC  health.UseCase
	cfg       *config.ServerConfig
	repo      repository.Repository
}

func NewManager(metricsUC metrics.UseCase, healthUC health.UseCase, cfg *config.ServerConfig, repo repository.Repository) *Manager {
	return &Manager{metricsUC: metricsUC, healthUC: healthUC, cfg: cfg, repo: repo}
}
