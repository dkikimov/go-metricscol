package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-metricscol/internal/proto"
	"go-metricscol/internal/server/health"
)

type HealthHandlers struct {
	healthUC health.UseCase
	proto.UnimplementedHealthServer
}

func NewHealthHandlers(healthUC health.UseCase) *HealthHandlers {
	return &HealthHandlers{healthUC: healthUC}
}

func (g HealthHandlers) Ping(ctx context.Context, _ *proto.PingRequest) (*proto.PingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := g.healthUC.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't ping db: %s", err)
	}

	return &proto.PingResponse{}, nil
}
