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
	if err != nil {
		apierror.WriteHTTP(w, err)
		return
	}

	metric, err := p.Storage.Get(urlData.MetricName, urlData.MetricType)
	if err != nil {
		apierror.WriteHTTP(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(metric.StringValue())); err != nil {
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

	foundMetric, err := p.Storage.Get(metric.Name, metric.MType)
	if err != nil {
		apierror.WriteHTTP(w, err)
		return
	}

	jsonFoundMetric, err := json.Marshal(foundMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonFoundMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Got metric with name %s, value: %s, type: %s", foundMetric.Name, foundMetric.StringValue(), foundMetric.MType)
}

func (p *Handlers) GetAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	getHashSubstring := func(metric models.Metric) string {
		if len(metric.Hash) == 0 {
			return ""
		}

		return fmt.Sprintf(", hash: %s", metric.Hash)
	}

	all, err := p.Storage.GetAll()
	if err != nil {
		apierror.WriteHTTP(w, err)
		return
	}

	for _, v := range all {
		_, err := w.Write([]byte(fmt.Sprintf("Key: %s, value: %s, type: %s%s \n", v.Name, v.StringValue(), v.MType, getHashSubstring(v))))
		if err != nil {
			log.Printf("Couldn't write response to GetAll request with error: %s", err)
		}
	}
}
