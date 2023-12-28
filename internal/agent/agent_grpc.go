package agent

import (
	"go-metricscol/internal/repository/memory"
)

type Grpc struct {
	cfg *Config
}

func NewGrpc(cfg *Config) *Grpc {
	return &Grpc{cfg: cfg}
}

func (agent Grpc) sendMetricsByOne(m *memory.Metrics) error {
	return nil
}

func (agent Grpc) sendMetricsAllTogether(m *memory.Metrics) error {
	return nil
}
