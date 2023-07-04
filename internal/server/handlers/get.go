package handlers

import (
	"fmt"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
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

func (p *Handlers) GetAll(w http.ResponseWriter, _ *http.Request) {
	for _, v := range p.Storage.GetAll() {
		_, err := w.Write([]byte(fmt.Sprintf("Key: %s, value: %s, type: %s \n", v.Name, v.GetStringValue(), v.MType)))
		if err != nil {
			log.Printf("Couldn't write response to GetAll request")
		}
	}
}
