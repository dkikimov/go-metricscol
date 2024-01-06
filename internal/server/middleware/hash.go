package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-metricscol/internal/models"
	"go-metricscol/internal/proto"
	"go-metricscol/internal/server/apierror"
)

// ValidateHashHandler is a middleware which gets models.Metric from request, calculates hash and compares it with given.
// If the hashes do not match, http.Error is called with code 400.
func (mw *Manager) ValidateHashHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashKey := mw.cfg.HashKey
		if len(hashKey) != 0 {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "couldn't read body", http.StatusInternalServerError)
				return
			}

			var metric models.Metric
			err = json.Unmarshal(body, &metric)
			if err != nil {
				http.Error(w, "couldn't parse json", http.StatusInternalServerError)
				return
			}

			if err := validateHashHandler(hashKey, metric); err != nil {
				apierror.WriteHTTP(w, err)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		next.ServeHTTP(w, r)
	}
}

func (mw *Manager) ValidateHashGrpcHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	updateRequest, err := req.(proto.UpdateRequest)
	if err {
		return handler(ctx, req)
	}

	hashKey := mw.cfg.HashKey
	if len(hashKey) != 0 {
		metric, err := proto.ParseMetricFromRequest(updateRequest.Metric)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "couldn't parse metric")
		}

		if err := validateHashHandler(hashKey, *metric); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Message)
		}
	}

	return handler(ctx, req)
}

func validateHashHandler(hashKey string, metric models.Metric) *apierror.APIError {
	if metric.HashValue(hashKey) != metric.Hash {
		return apierror.NewAPIError(http.StatusBadRequest, "hash mismatch")
	}

	return nil
}

// ValidateHashesHandler is a middleware which gets []models.Metric from request, calculates hashes and compares them with given.
// If at least one of the hashes do not match, http.Error is called with code 400.
func (mw *Manager) ValidateHashesHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashKey := mw.cfg.HashKey
		if len(hashKey) != 0 {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "couldn't read body", http.StatusInternalServerError)
				return
			}

			var metrics []models.Metric
			err = json.Unmarshal(body, &metrics)
			if err != nil {
				http.Error(w, "couldn't parse json", http.StatusInternalServerError)
				return
			}

			if err := validateHashesHandler(hashKey, metrics); err != nil {
				apierror.WriteHTTP(w, err)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		next.ServeHTTP(w, r)
	}
}

func (mw *Manager) ValidateHashesGrpcHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	updateRequest, err := req.(proto.UpdatesRequest)
	if err {
		return handler(ctx, req)
	}

	hashKey := mw.cfg.HashKey
	if len(hashKey) != 0 {
		var metrics []models.Metric

		for _, metric := range updateRequest.Metric {
			m, err := proto.ParseMetricFromRequest(metric)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "couldn't parse metrics")
			}
			metrics = append(metrics, *m)
		}

		if err := validateHashesHandler(hashKey, metrics); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Message)
		}
	}

	return handler(ctx, req)
}

func validateHashesHandler(hashKey string, metrics []models.Metric) *apierror.APIError {
	for _, metric := range metrics {
		if metric.HashValue(hashKey) != metric.Hash {
			return apierror.NewAPIError(http.StatusBadRequest, "hash mismatch")
		}
	}

	return nil
}
