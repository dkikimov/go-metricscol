package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-metricscol/internal/config"
	"go-metricscol/internal/models"
	"go-metricscol/internal/proto"
	"go-metricscol/internal/server/metrics"
)

type MetricsHandlers struct {
	metricsUC metrics.UseCase
	config    *config.ServerConfig
	proto.UnimplementedMetricsServer
}

func (g MetricsHandlers) UpdateMetric(ctx context.Context, request *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	var response proto.UpdateResponse

	requestMetric, err := proto.ParseMetricFromRequest(request.Metric)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't parse metric from request: %s", err)
	}

	if err := g.metricsUC.Update(ctx, *requestMetric); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't update metric: %s", err)
	}

	return &response, nil
}

func (g MetricsHandlers) UpdatesMetric(ctx context.Context, request *proto.UpdatesRequest) (*proto.UpdatesResponse, error) {
	var response proto.UpdatesResponse

	var requestMetrics = make([]models.Metric, len(request.Metric))
	for i, metric := range request.Metric {
		requestMetric, err := proto.ParseMetricFromRequest(metric)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "couldn't parse metric from request: %s", err)
		}

		requestMetrics[i] = *requestMetric
	}

	if err := g.metricsUC.Updates(ctx, requestMetrics); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't update metric: %w", err)
	}

	return &response, nil
}

func (g MetricsHandlers) ValueMetric(ctx context.Context, request *proto.ValueRequest) (*proto.ValueResponse, error) {
	var response proto.ValueResponse

	metricType, err := proto.ParseTypeFromRequest(request.Type)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "couldn't parse metric type from request: %s", err)
	}

	foundMetric, err := g.metricsUC.Find(ctx, request.Name, metricType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't find metric: %s", err)
	}

	response.Metric = &proto.Metric{
		Name:  foundMetric.Name,
		Type:  proto.MetricType(foundMetric.MType.IntGrpc()),
		Value: foundMetric.StringValue(),
		Hash:  foundMetric.HashValue(g.config.HashKey),
	}

	return &response, nil
}

func (g MetricsHandlers) ListMetrics(context.Context, *proto.ListRequest) (*proto.ListResponse, error) {
	var response proto.ListResponse

	metricsList, err := g.metricsUC.GetAll(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't list metrics: %s", err)
	}

	response.Metric = make([]*proto.Metric, len(metricsList))
	for i, metric := range metricsList {
		response.Metric[i] = &proto.Metric{
			Name:  metric.Name,
			Type:  proto.MetricType(metric.MType.IntGrpc()),
			Value: metric.StringValue(),
			Hash:  metric.HashValue(g.config.HashKey),
		}
	}

	return &response, nil
}

func NewMetricsHandlers(metricsUC metrics.UseCase, config *config.ServerConfig) *MetricsHandlers {
	return &MetricsHandlers{metricsUC: metricsUC, config: config}
}
