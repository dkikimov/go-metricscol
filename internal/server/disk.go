package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

func (s Server) enableSavingToDisk(ctx context.Context) error {
	if !s.Repository.SupportsSavingToDisk() {
		return errors.New("selected repository doesn't support saving to disk")
	}

	ticker := time.NewTicker(s.Config.StoreInterval)

	for {
		select {
		case <-ticker.C:
			if err := s.saveToDisk(); err != nil {
				log.Printf("Couldn't save metrics to disk with error: %s", err)
			}

		case <-ctx.Done():
			if err := s.saveToDisk(); err != nil {
				log.Printf("Couldn't save metrics to disk with error: %s", err)
			}
			return nil
		}
	}
}

func (s Server) saveToDisk() error {
	if !s.Repository.SupportsSavingToDisk() {
		log.Printf("Selected repository doesn't support saving to disk")
		return errors.New("unsupported storage")
	}

	log.Printf("saving to disk")

	file, err := os.OpenFile(s.Config.StoreFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(s.Repository); err != nil {
		return err
	}

	return nil
}

func (s Server) restoreFromDisk() error {
	file, err := os.OpenFile(s.Config.StoreFile, os.O_RDONLY|os.O_SYNC, 0777)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(s.Repository); err != nil {
		return err
	}

	return nil
}

func (s Server) diskSaverHandler(next http.HandlerFunc, _ *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		saveToDisk := s.Config.StoreInterval == 0 && len(s.Config.StoreFile) != 0 && len(s.Config.DatabaseDSN) == 0
		if saveToDisk {
			if err := s.saveToDisk(); err != nil {
				log.Printf("Couldn't save metrics to disk with error: %s", err)
			}
		}
	}
}
