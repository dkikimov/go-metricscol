package handlers

import (
	"context"
	"net/http"
	"time"
)

func (p *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := p.Postgres.Ping(ctx)
	if err != nil {
		http.Error(w, "couldn't ping db", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
