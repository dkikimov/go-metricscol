package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"go-metricscol/internal/server/health"
)

type HealthHandlers struct {
	healthUC health.UseCase
}

func NewHealthHandlers(healthUC health.UseCase) *HealthHandlers {
	return &HealthHandlers{healthUC: healthUC}
}

func (h HealthHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err := h.healthUC.Ping(ctx)
	if err != nil {
		http.Error(w, "couldn't ping db", http.StatusInternalServerError)
		log.Printf("Couldn't ping db with error: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
