package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
)

// Update is a handler that updates models.Metric with given key based on the parameters in the URL.
func (p *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	urlData, err := models.ParsePostURLData(r)
	if err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't parse url with error: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err = p.Storage.Update(ctx, urlData.MetricName, urlData.MetricType, urlData.MetricValue); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Update metric with name %s, value: %s, type: %s", urlData.MetricName, urlData.MetricValue, urlData.MetricType)
	w.WriteHeader(http.StatusOK)
}

// UpdateJSON is a handler that updates models.Metric with the given key based on the json in the request body.
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := p.Storage.UpdateWithStruct(ctx, &metric); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Update metric with name %s, value: %s, type: %s", metric.Name, metric.StringValue(), metric.MType)

	ctxGet, cancelGet := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancelGet()

	newMetric, _ := p.Storage.Get(ctxGet, metric.Name, metric.MType)

	p.addHash(newMetric)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newMetric)
	if err != nil {
		http.Error(w, "couldn't encode json", http.StatusInternalServerError)
		log.Printf("Couldn't encode json with error: %s", err)
		return
	}
}

// Updates is a handler that updates []models.Metric based on the json in the request body.
func (p *Handlers) Updates(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read body", http.StatusInternalServerError)
		log.Printf("Couldn't read body with error: %s", err)
		return
	}

	var metrics []models.Metric
	if err := json.Unmarshal(body, &metrics); err != nil {
		http.Error(w, "couldn't parse json", http.StatusBadRequest)
		log.Printf("Couldn't parse json with error: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := p.Storage.Updates(ctx, metrics); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Updates %d metrics", len(metrics))

	updatedMetrics := make([]models.Metric, 0, len(metrics))

	ctxGet, cancelGet := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancelGet()

	for _, metric := range metrics {
		newMetric, _ := p.Storage.Get(ctxGet, metric.Name, metric.MType)

		p.addHash(newMetric)
		updatedMetrics = append(updatedMetrics, *newMetric)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedMetrics)
	if err != nil {
		http.Error(w, "couldn't encode json", http.StatusInternalServerError)
		log.Printf("Couldn't encode json with error: %s", err)
		return
	}
}
