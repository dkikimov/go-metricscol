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
		log.Printf("Couldn't parse url with error: %s", err)
		return
	}

	if err = p.Storage.Update(urlData.MetricName, urlData.MetricType, urlData.MetricValue); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Update metric with name %s, value: %s, type: %s", urlData.MetricName, urlData.MetricValue, urlData.MetricType)
	w.WriteHeader(http.StatusOK)
}

func (p *Handlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read body", http.StatusInternalServerError)
		log.Printf("Couldn't read body with error: %s", err)
		return
	}

	var metric models.Metric
	if err := json.Unmarshal(body, &metric); err != nil {
		http.Error(w, "couldn't parse json", http.StatusBadRequest)
		log.Printf("Couldn't parse json with error: %s", err)
		return
	}

	if err := p.Storage.UpdateWithStruct(&metric); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Update metric with name %s, value: %s, type: %s", metric.Name, metric.StringValue(), metric.MType)

	newMetric, _ := p.Storage.Get(metric.Name, metric.MType)

	err = json.NewEncoder(w).Encode(newMetric)
	if err != nil {
		http.Error(w, "couldn't encode json", http.StatusInternalServerError)
		log.Printf("Couldn't encode json with error: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
