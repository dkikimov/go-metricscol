package server

import (
	"context"
	"errors"
	"log"
	"time"
)

func (s Server) enableSavingToDisk(ctx context.Context) error {
	if !s.Repo.SupportsSavingToDisk() {
		return errors.New("selected repository doesn't support saving to disk")
	}

	ticker := time.NewTicker(s.Config.StoreInterval)

	for {
		select {
		case <-ticker.C:
			if err := s.Repo.SaveToDisk(s.Config.StoreFile); err != nil {
				log.Printf("Couldn't save metrics to disk with error: %s", err)
			}

		case <-ctx.Done():
			if err := s.Repo.SaveToDisk(s.Config.StoreFile); err != nil {
				log.Printf("Couldn't save metrics to disk with error: %s", err)
			}
			return nil
		}
	}
}
