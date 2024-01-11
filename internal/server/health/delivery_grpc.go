package health

import "go-metricscol/internal/proto"

type GrpcHandlers interface {
	proto.HealthServer
}
