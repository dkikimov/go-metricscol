package server

import (
	"go-metricscol/internal/repository"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Config     *Config
	Repository repository.Repository
}

func NewServer(config *Config, repository repository.Repository) *Server {
	return &Server{Config: config, Repository: repository}
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
