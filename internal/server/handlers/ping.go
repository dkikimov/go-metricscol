package handlers

import (
	"context"
	"log"
	"net/http"
	"time"
)

func (p *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err := p.Postgres.Ping(ctx)
	if err != nil {
		http.Error(w, "couldn't ping db", http.StatusInternalServerError)
		log.Printf("Couldn't ping db with error: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
