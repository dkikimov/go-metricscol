package server

import (
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/repository/postgres"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Config     *Config
	Repository repository.Repository
	Postgres   *postgres.DB
}

func NewServer(config *Config) (*Server, error) {
	db, err := postgres.New(config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &Server{Config: config, Repository: getRepository(config, db), Postgres: db}, nil
}

func getRepository(config *Config, db *postgres.DB) repository.Repository {
	if len(config.DatabaseDSN) == 0 {
		return memory.NewMemStorage(config.HashKey)
	} else {
		return db
	}
}

func (s Server) ListenAndServe() error {
	r := s.newRouter(s.Repository)

	serv := http.Server{
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

	if len(s.Config.StoreFile) != 0 && s.Config.StoreInterval != 0 && len(s.Config.DatabaseDSN) == 0 {
		go s.enableSavingToDisk()
	}

	return serv.ListenAndServe()
}
