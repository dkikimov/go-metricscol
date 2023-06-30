package handlers

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/server/apierror"
	"log"
	"net/http"
)

type Processors struct {
	Storage repository.Repository
}

func (p *Processors) Update(w http.ResponseWriter, r *http.Request) {
	urlData, err := models.ParseURLData(r)
	if err != apierror.NoError {
		w.WriteHeader(err.StatusCode())
		return
	}

	if err := p.Storage.Update(urlData.MetricName, urlData.MetricValue, urlData.MetricType); err != apierror.NoError {
		w.WriteHeader(err.StatusCode())
		return
	}

	log.Printf("Updated metric with name %s, value: %s, type: %s", urlData.MetricName, urlData.MetricValue, urlData.MetricType)
	w.WriteHeader(http.StatusOK)
}

func (p *Processors) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
