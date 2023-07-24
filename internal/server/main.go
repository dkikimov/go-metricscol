package server

import (
	"github.com/jackc/pgx/v5"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/postgres"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Config     *Config
	Repository repository.Repository
	Postgres   *pgx.Conn
}

func NewServer(config *Config, repository repository.Repository) (*Server, error) {
	db, err := postgres.New(config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &Server{Config: config, Repository: repository, Postgres: db}, nil
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

	if len(s.Config.StoreFile) != 0 && s.Config.StoreInterval != 0 {
		go s.enableSavingToDisk()
	}

	return serv.ListenAndServe()
}
