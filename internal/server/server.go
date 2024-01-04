package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
	"golang.org/x/sync/errgroup"

	"go-metricscol/internal/config"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/postgres"
)

// Server defines config and repository for HTTP server instance.
type Server struct {
	Config   *config.ServerConfig
	Repo     repository.Repository
	Postgres *postgres.DB
}

// NewServer returns new Server with defined config.
func NewServer(config *config.ServerConfig, repo repository.Repository, postgres *postgres.DB) *Server {
	return &Server{Config: config, Repo: repo, Postgres: postgres}
}

// func createBackendBasedOnType(cfg *repository.Repository, backendType BackendType) (Backend, error) {
// 	switch backendType {
// 	case GRPC:
// 		return NewGrpc(cfg)
// 	case HTTP:
// 		return newHttpRouter(cfg), nil
// 	default:
// 		return nil, fmt.Errorf("unknown backend type id: %d", backendType)
// 	}
// }

// ListenAndServe listens on the TCP network address given in config and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (s Server) ListenAndServe(ctx context.Context) error {
	r := chi.NewRouter()

	if err := s.MapHandlers(r); err != nil {
		return errors.New("couldn't map handlers")
	}

	httpServer := http.Server{
		Addr:    s.Config.Address,
		Handler: r,
	}

	if s.Config.Restore {
		if err := s.Repo.RestoreFromDisk(s.Config.StoreFile); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("error while restoring from disk: %s", err)
			}
		}
	}

	group, _ := errgroup.WithContext(context.Background())

	shutdownWg := sync.WaitGroup{}
	diskContext, cancel := context.WithCancel(context.Background())
	if len(s.Config.StoreFile) != 0 && s.Config.StoreInterval != 0 && len(s.Config.DatabaseDSN) == 0 {
		group.Go(func() error {
			shutdownWg.Add(1)
			defer shutdownWg.Done()

			err := s.enableSavingToDisk(diskContext)
			if err != nil {
				return fmt.Errorf("couldn't enable saving to disk: %s", err)
			}
			return nil
		})
	}

	group.Go(func() error {
		if err := httpServer.ListenAndServe(); err != nil {
			return err
		}

		return nil
	})

	<-ctx.Done()
	err := httpServer.Shutdown(context.Background())
	cancel()

	shutdownWg.Wait()
	if err := group.Wait(); err != nil {
		return err
	}

	return err
}
