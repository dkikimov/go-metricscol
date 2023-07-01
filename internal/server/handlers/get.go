package handlers

import (
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"log"
	"net/http"
)

func (p *Processors) Get(w http.ResponseWriter, r *http.Request) {
	urlData, apiError := models.ParseGetURLData(r)
	if apiError != apierror.NoError {
		w.WriteHeader(apiError.StatusCode())
		return
	}

	metricValue, apiError := p.Storage.GetString(urlData.MetricName, urlData.MetricType)
	if apiError != apierror.NoError {
		w.WriteHeader(apiError.StatusCode())
		return
	}

	if _, err := w.Write([]byte(metricValue)); err != nil {
		log.Printf("Couldn't write response")
	}
}
