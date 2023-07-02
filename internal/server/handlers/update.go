package handlers

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"log"
	"net/http"
)

func (p *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	urlData, err := models.ParsePostURLData(r)
	if err != apierror.NoError {
		w.WriteHeader(err.StatusCode())
		return
	}

	if err := p.Storage.Update(urlData.MetricName, urlData.MetricType, urlData.MetricValue); err != apierror.NoError {
		w.WriteHeader(err.StatusCode())
		return
	}

	log.Printf("Updated metric with name %s, value: %s, type: %s", urlData.MetricName, urlData.MetricValue, urlData.MetricType)
	w.WriteHeader(http.StatusOK)
}
