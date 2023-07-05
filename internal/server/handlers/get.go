package handlers

import (
	"encoding/json"
	"fmt"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"io"
	"log"
	"net/http"
)

func (p *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	urlData, err := models.ParseGetURLData(r)
	if err != apierror.NoError {
		w.WriteHeader(err.StatusCode())
		return
	}

	metric, err := p.Storage.Get(urlData.MetricName, urlData.MetricType)
	if err != apierror.NoError {
		w.WriteHeader(err.StatusCode())
		return
	}

	if _, err := w.Write([]byte(metric.GetStringValue())); err != nil {
		log.Printf("Couldn't write response")
	}
}

func (p *Handlers) GetJSON(w http.ResponseWriter, r *http.Request) {
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

	foundMetric, apiError := p.Storage.Get(metric.Name, metric.MType)
	if apiError != apierror.NoError {
		http.Error(w, "Not found", apiError.StatusCode())
		return
	}

	jsonFoundMetric, err := json.Marshal(foundMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonFoundMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Got metric with name %s, value: %s, type: %s", foundMetric.Name, foundMetric.GetStringValue(), foundMetric.MType)
}

func (p *Handlers) GetAll(w http.ResponseWriter, _ *http.Request) {
	for _, v := range p.Storage.GetAll() {
		_, err := w.Write([]byte(fmt.Sprintf("Key: %s, value: %s, type: %s \n", v.Name, v.GetStringValue(), v.MType)))
		if err != nil {
			log.Printf("Couldn't write response to GetAll request")
		}
	}
}
