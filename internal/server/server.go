package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"go-metricscol/internal/config"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/backends"
)

// Server defines config and repository for HTTPType server instance.
type Server struct {
	Config  *config.ServerConfig
	Repo    repository.Repository
	Backend backends.Backend
}

// NewServer returns new Server with defined config.
func NewServer(config *config.ServerConfig, repo repository.Repository, backendType backends.Backend) *Server {
	return &Server{Config: config, Repo: repo, Backend: backendType}
}

// ListenAndServe listens on the TCP network address given in config and then calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (s Server) ListenAndServe(ctx context.Context) error {
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
		if err := s.Backend.ListenAndServe(); err != nil {
			return err
		}

		return nil
	})

	<-ctx.Done()
	err := s.Backend.GracefulShutdown(context.Background())
	cancel()

	shutdownWg.Wait()
	if err := group.Wait(); err != nil {
		return err
	}

	return err
}
