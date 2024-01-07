package http

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-metricscol/internal/config"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"go-metricscol/internal/server/metrics"
)

type MetricsHandlers struct {
	metricsUC metrics.UseCase
	config    *config.ServerConfig
}

func NewMetricsHandlers(metricsUC metrics.UseCase, config *config.ServerConfig) *MetricsHandlers {
	return &MetricsHandlers{metricsUC: metricsUC, config: config}
}

// Find is a handler that finds models.Metric based on the parameters in the URL.
// If metric is not found 404 status code returned.
// Otherwise, metric value is returned.
func (m *MetricsHandlers) Find(w http.ResponseWriter, r *http.Request) {
	urlData, err := models.ParseGetURLData(r)
	if err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't parse url with error: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	metric, err := m.metricsUC.Find(ctx, urlData.MetricName, urlData.MetricType)
	if err != nil {
		apierror.WriteHTTP(w, err)
		if !errors.Is(err, apierror.NotFound) {
			log.Printf("Couldn't get metric with error: %s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(metric.StringValue())); err != nil {
		log.Printf("Couldn't write response")
	}
}

// FindJSON is a handler that finds models.Metric based on the json passed in request body.
// In case the json could not be parsed, the status code 400 is returned.
// If metric is not found 404 status code returned.
// Otherwise, metric value is returned.
func (m *MetricsHandlers) FindJSON(w http.ResponseWriter, r *http.Request) {
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

	foundMetric, err := m.metricsUC.Find(ctx, metric.Name, metric.MType)
	if err != nil {
		apierror.WriteHTTP(w, err)
		if !errors.Is(err, apierror.NotFound) {
			log.Printf("Couldn't get metric with error: %s", err)
		}
		return
	}

	foundMetric.Hash = foundMetric.HashValue(m.config.HashKey)

	jsonFoundMetric, err := json.Marshal(foundMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Couldn't marshal metric with error: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonFoundMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Couldn't write response to FindJSON request with error: %s", err)
		return
	}

	log.Printf("Got metric with name %s, value: %s, type: %s", foundMetric.Name, foundMetric.StringValue(), foundMetric.MType)
}

// Update is a handler that updates models.Metric with given key based on the parameters in the URL.
func (m *MetricsHandlers) Update(w http.ResponseWriter, r *http.Request) {
	urlData, err := models.ParsePostURLData(r)
	if err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't parse url with error: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	metric := models.Metric{
		Name:  urlData.MetricName,
		MType: urlData.MetricType,
	}

	switch urlData.MetricType {
	case models.Gauge:
		floatVal, err := strconv.ParseFloat(urlData.MetricValue, 64)
		if err != nil {
			apierror.WriteHTTP(w, apierror.NumberParse)
			return
		}

		metric.Value = &floatVal
	case models.Counter:
		intVal, err := strconv.ParseInt(urlData.MetricValue, 10, 64)
		if err != nil {
			apierror.WriteHTTP(w, apierror.NumberParse)
			return
		}

		metric.Delta = &intVal
	default:
		apierror.WriteHTTP(w, apierror.UnknownMetricType)
		return
	}

	if err = m.metricsUC.Update(ctx, metric); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Update metric with name %s, value: %s, type: %s", urlData.MetricName, urlData.MetricValue, urlData.MetricType)
	w.WriteHeader(http.StatusOK)
}

// UpdateJSON is a handler that updates models.Metric with the given key based on the json in the request body.
func (m *MetricsHandlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "couldn't read body", http.StatusInternalServerError)
		log.Printf("Couldn't read body with error: %s", err)
		return
	}

	var decryptedJSON []byte
	if m.config.CryptoKey != nil {
		decryptedJSON, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, m.config.CryptoKey, body, nil)
		if err != nil {
			http.Error(w, "couldn't decrypt json", http.StatusInternalServerError)
			log.Printf("couldn't decrypt json: %s", err)
			return
		}
	} else {
		decryptedJSON = body
	}

	var metric models.Metric
	if err := json.Unmarshal(decryptedJSON, &metric); err != nil {
		http.Error(w, "couldn't parse json", http.StatusBadRequest)
		log.Printf("Couldn't parse json with error: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := m.metricsUC.Update(ctx, metric); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Update metric with name %s, value: %s, type: %s", metric.Name, metric.StringValue(), metric.MType)

	ctxGet, cancelGet := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancelGet()

	newMetric, _ := m.metricsUC.Find(ctxGet, metric.Name, metric.MType)
	newMetric.Hash = newMetric.HashValue(m.config.HashKey)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newMetric)
	if err != nil {
		http.Error(w, "couldn't encode json", http.StatusInternalServerError)
		log.Printf("Couldn't encode json with error: %s", err)
		return
	}
}

// Updates is a handler that updates []models.Metric based on the json in the request body.
func (m *MetricsHandlers) Updates(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read body", http.StatusInternalServerError)
		log.Printf("Couldn't read body with error: %s", err)
		return
	}

	var metricSlice []models.Metric
	if err := json.Unmarshal(body, &metricSlice); err != nil {
		http.Error(w, "couldn't parse json", http.StatusBadRequest)
		log.Printf("Couldn't parse json with error: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := m.metricsUC.Updates(ctx, metricSlice); err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't update metric with error: %s", err)
		return
	}

	log.Printf("Updates %d metrics", len(metricSlice))

	updatedMetrics := make([]models.Metric, 0, len(metricSlice))

	ctxGet, cancelGet := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancelGet()

	for _, metric := range metricSlice {
		newMetric, _ := m.metricsUC.Find(ctxGet, metric.Name, metric.MType)
		newMetric.Hash = newMetric.HashValue(m.config.HashKey)
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

// GetAll returns list of all metrics stored on repository.
func (m *MetricsHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	getHashSubstring := func(metric models.Metric) string {
		if len(metric.Hash) == 0 {
			return ""
		}

		return fmt.Sprintf(", hash: %s", metric.Hash)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	all, err := m.metricsUC.GetAll(ctx)
	if err != nil {
		apierror.WriteHTTP(w, err)
		log.Printf("Couldn't get all metrics with error: %s", err)
		return
	}

	for idx, value := range all {
		all[idx].Hash = value.HashValue(m.config.HashKey)
	}

	for _, v := range all {
		_, err := w.Write([]byte(fmt.Sprintf("Key: %s, value: %s, type: %s%s \n", v.Name, v.StringValue(), v.MType, getHashSubstring(v))))
		if err != nil {
			log.Printf("Couldn't write response to GetAll request with error: %s", err)
		}
	}
}
