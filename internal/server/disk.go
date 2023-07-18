package server

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func (s Server) enableSavingToDisk() {
	ticker := time.NewTicker(s.Config.StoreInterval)

	for range ticker.C {
		if err := s.saveToDisk(); err != nil {
			log.Printf("Couldn't save metrics to disk with error: %s", err)
		}
	}
}

func (s Server) saveToDisk() error {
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
