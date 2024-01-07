package backends

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"go-metricscol/internal/config"
	"go-metricscol/internal/proto"
	"go-metricscol/internal/repository"
	healthGrpc "go-metricscol/internal/server/health/delivery/grpc"
	helathUseCase "go-metricscol/internal/server/health/usecase"
	metricsGrpc "go-metricscol/internal/server/metrics/delivery/grpc"
	metricsUseCase "go-metricscol/internal/server/metrics/usecase"
	"go-metricscol/internal/server/middleware"
)

type Grpc struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGrpc(repo repository.Repository, config *config.ServerConfig, listener net.Listener) (*Grpc, error) {
	metricsUC := metricsUseCase.NewMetricsUC(repo, config)
	healthUC := helathUseCase.NewHealthUC(repo)

	mw := middleware.NewManager(metricsUC, healthUC, config, repo)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.DiskSaverGrpcMiddleware,
			mw.ValidateHashGrpcHandler,
			mw.GrpcTrustedSubnetHandler,
			mw.ValidateHashesGrpcHandler,
		))

	proto.RegisterHealthServer(server, healthGrpc.NewHealthHandlers(healthUC))
	proto.RegisterMetricsServer(server, metricsGrpc.NewMetricsHandlers(metricsUC, config))

	return &Grpc{server: server, listener: listener}, nil
}

func (s Grpc) ListenAndServe() error {
	return s.server.Serve(s.listener)
}

func (s Grpc) GracefulShutdown(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}
