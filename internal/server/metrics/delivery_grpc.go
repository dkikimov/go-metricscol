package metrics

import "go-metricscol/internal/proto"

type GrpcHandlers interface {
	proto.MetricsServer
}
