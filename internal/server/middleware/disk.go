package middleware

import (
	"log"
	"net/http"
)

func (mw *Manager) DiskSaverHttpMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		saveToDisk := mw.cfg.StoreInterval == 0 && len(mw.cfg.StoreFile) != 0 && len(mw.cfg.DatabaseDSN) == 0
		if saveToDisk {
			if err := mw.repo.SaveToDisk(mw.cfg.StoreFile); err != nil {
				log.Printf("Couldn't save metrics to disk with error: %s", err)
			}
		}
	}
}
