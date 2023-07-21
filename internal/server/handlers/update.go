package handlers

import (
	"encoding/json"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"io"
	"log"
	"net/http"
)

func (p *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	urlData, err := models.ParsePostURLData(r)
	if err != nil {
		apierror.WriteHTTP(w, err)
		return
	}

	if err = p.Storage.Update(urlData.MetricName, urlData.MetricType, urlData.MetricValue); err != nil {
		apierror.WriteHTTP(w, err)
		return
	}

	log.Printf("Updated metric with name %s, value: %s, type: %s", urlData.MetricName, urlData.MetricValue, urlData.MetricType)
	w.WriteHeader(http.StatusOK)
}

func (p *Handlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read body", http.StatusInternalServerError)
		return
	}

	var metric models.Metric
	if err := json.Unmarshal(body, &metric); err != nil {
		http.Error(w, "couldn't parse json", http.StatusBadRequest)
		return
	}

	if err := p.Storage.UpdateWithStruct(&metric); err != nil {
		apierror.WriteHTTP(w, err)
		return
	}

	newMetric, _ := p.Storage.Get(metric.Name, metric.MType)
	log.Printf("Updated metric with name %s, value: %s, type: %s", newMetric.Name, newMetric.StringValue(), newMetric.MType)
	w.WriteHeader(http.StatusOK)
}
