package agent

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "go-metricscol/internal/proto"
	"go-metricscol/internal/repository/memory"
)

type Grpc struct {
	cfg    *Config
	conn   *grpc.ClientConn
	client pb.MetricsClient
}

func NewGrpc(cfg *Config) (*Grpc, error) {
	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Grpc{cfg: cfg, conn: conn, client: pb.NewMetricsClient(conn)}, nil
}

func (agent Grpc) SendMetricsByOne(m *memory.Metrics) error {
	for _, value := range m.Collection {
		metric := pb.Metric{
			Name:  value.Name,
			Type:  pb.MetricType(value.MType.IntGrpc()),
			Value: value.StringValue(),
			Hash:  value.HashValue(agent.cfg.HashKey),
		}

		_, err := agent.client.UpdateMetric(context.Background(), &pb.UpdateRequest{Metric: &metric})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				if e.Code() != codes.OK {
					return fmt.Errorf("coudln't send metrics, status code: %d, response: %s", e.Code(), e.Message())
				}
			}
		}
	}

	return nil
}

func (agent Grpc) SendMetricsAllTogether(m *memory.Metrics) error {
	metrics := make([]*pb.Metric, 0, len(m.Collection))
	for _, value := range m.Collection {
		metric := pb.Metric{
			Name:  value.Name,
			Type:  pb.MetricType(value.MType.IntGrpc()),
			Value: value.StringValue(),
			Hash:  value.HashValue(agent.cfg.HashKey),
		}

		metrics = append(metrics, &metric)
	}

	_, err := agent.client.UpdatesMetric(context.Background(), &pb.UpdatesRequest{Metric: metrics})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() != codes.OK {
				return fmt.Errorf("coudln't send metrics, status code: %d, response: %s", e.Code(), e.Message())
			}
		}
	}

	return nil
}

func (agent Grpc) Close() error {
	return agent.conn.Close()
}
