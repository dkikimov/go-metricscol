package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/repository/postgres"
)

// Server defines config and repository for HTTP server instance.
type Server struct {
	Config     *Config
	Repository repository.Repository
	Postgres   *postgres.DB
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

// NewServer returns new Server with defined config.
// Initialized database if necessary.
func NewServer(config *Config) (*Server, error) {
	db, err := postgres.New(config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &Server{Config: config, Repository: getRepository(config, db), Postgres: db}, nil
}

func getRepository(config *Config, db *postgres.DB) repository.Repository {
	if len(config.DatabaseDSN) == 0 {
		return memory.NewMemStorage()
	} else {
		return db
	}
}

// ListenAndServe listens on the TCP network address given in config and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (s Server) ListenAndServe(ctx context.Context) error {
	r := s.newHttpRouter(s.Repository)

	httpServer := http.Server{
		Addr:    s.Config.Address,
		Handler: r,
	}

	if s.Config.Restore {
		if err := s.restoreFromDisk(); err != nil {
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
