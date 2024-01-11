package agent

import "go-metricscol/internal/repository/memory"

type BackendType int

const (
	GRPC BackendType = iota
	HTTP
)

type Backend interface {
	SendMetricsByOne(m *memory.Metrics) error
	SendMetricsAllTogether(m *memory.Metrics) error
	Close() error
}
