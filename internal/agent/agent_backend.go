package agent

import "go-metricscol/internal/repository/memory"

type Backend interface {
	sendMetricsByOne(m *memory.Metrics) error
	sendMetricsAllTogether(m *memory.Metrics) error
}
