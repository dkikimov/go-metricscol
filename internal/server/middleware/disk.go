package middleware

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-metricscol/internal/config"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/apierror"
)

func (mw *Manager) DiskSaverHttpMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if err := diskSaverMiddleware(mw.cfg, mw.repo); err != nil {
			apierror.WriteHTTP(w, err)
			return
		}
	}
}

func (mw *Manager) DiskSaverGrpcMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)

	if err := diskSaverMiddleware(mw.cfg, mw.repo); err != nil {
		// TODO: implement own code errors system
		return nil, status.Errorf(codes.Internal, err.Message)
	}

	return resp, err
}

func diskSaverMiddleware(cfg *config.ServerConfig, repository repository.Repository) *apierror.APIError {
	saveToDisk := cfg.StoreInterval == 0 && len(cfg.StoreFile) != 0 && len(cfg.DatabaseDSN) == 0
	if saveToDisk {
		if err := repository.SaveToDisk(cfg.StoreFile); err != nil {
			return apierror.NewAPIError(http.StatusInternalServerError, fmt.Sprintf("Couldn't save metrics to disk with error: %s", err))
		}
	}

	return nil
}
