package grpcpackage

import (
	"context"
	"fmt"

	"go-metricscol/internal/config"
	"go-metricscol/internal/models"
	"go-metricscol/internal/proto"
	"go-metricscol/internal/server/metrics"
)

type GrpcMetricsHandlers struct {
	metricsUC metrics.UseCase
	config    *config.ServerConfig
	proto.UnimplementedMetricsServer
}

func (g GrpcMetricsHandlers) UpdateMetric(ctx context.Context, request *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	var response proto.UpdateResponse

	requestMetric, err := parseMetricFromRequest(request.Metric)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse metric from request: %w", err)
	}

	if err := g.metricsUC.Update(ctx, *requestMetric); err != nil {
		return nil, fmt.Errorf("couldn't update metric: %w", err)
	}

	return &response, nil
}

func (g GrpcMetricsHandlers) UpdatesMetric(ctx context.Context, request *proto.UpdatesRequest) (*proto.UpdatesResponse, error) {
	var response proto.UpdatesResponse

	var requestMetrics = make([]models.Metric, len(request.Metric))
	for i, metric := range request.Metric {
		requestMetric, err := parseMetricFromRequest(metric)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse metric from request: %w", err)
		}

		requestMetrics[i] = *requestMetric
	}

	if err := g.metricsUC.Updates(ctx, requestMetrics); err != nil {
		return nil, fmt.Errorf("couldn't update metric: %w", err)
	}

	return &response, nil
}

func (g GrpcMetricsHandlers) ValueMetric(ctx context.Context, request *proto.ValueRequest) (*proto.ValueResponse, error) {
	var response proto.ValueResponse

	metricType, err := parseTypeFromRequest(request.Type)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse metric type from request: %w", err)
	}

	foundMetric, err := g.metricsUC.Find(ctx, request.Name, metricType)
	if err != nil {
		return nil, fmt.Errorf("couldn't find metric: %w", err)
	}

	response.Metric = &proto.Metric{
		Name:  foundMetric.Name,
		Type:  proto.MetricType(foundMetric.MType.IntGrpc()),
		Value: foundMetric.StringValue(),
		Hash:  foundMetric.HashValue(g.config.HashKey),
	}

	return &response, nil
}

func NewGrpcMetricsHandlers(metricsUC metrics.UseCase, config *config.ServerConfig) *GrpcMetricsHandlers {
	return &GrpcMetricsHandlers{metricsUC: metricsUC, config: config}
}
