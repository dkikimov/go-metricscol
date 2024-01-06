package middleware

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go-metricscol/internal/config"
	"go-metricscol/internal/server/apierror"
)

func (mw *Manager) HttpTrustedSubnetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get("X-Real-IP")

		if err := trustedSubnetHandler(mw.cfg, headerValue); err != nil {
			apierror.WriteHTTP(w, err)
			return
		}

		r.Header.Del("X-Real-IP")

		next.ServeHTTP(w, r)
	})
}

func (mw *Manager) GrpcTrustedSubnetHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var headerValue string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			headerValue = values[0]
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "couldn't parse context")
	}

	if err := trustedSubnetHandler(mw.cfg, headerValue); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Message)
	}

	return handler(ctx, req)
}

func trustedSubnetHandler(cfg *config.ServerConfig, headerValue string) *apierror.APIError {
	if len(headerValue) != 0 && len(cfg.TrustedSubnet) != 0 {
		ip := net.ParseIP(headerValue)
		if ip == nil {
			return apierror.NewAPIError(http.StatusBadRequest, "Couldn't parse IP from header X-Real-IP")
		}

		_, ipNet, err := net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			return apierror.NewAPIError(http.StatusInternalServerError, "Couldn't parse CIDR")
		}

		if !ipNet.Contains(ip) {
			return apierror.NewAPIError(http.StatusForbidden, "X-Real-IP is not in trusted subnet")
		}
	}
	return nil
}
