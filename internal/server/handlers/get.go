package handlers

import (
	"fmt"
	"go-metricscol/internal/models"
	"go-metricscol/internal/server/apierror"
	"log"
	"net/http"
	"sort"
)

func (p *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	urlData, apiError := models.ParseGetURLData(r)
	if apiError != apierror.NoError {
		w.WriteHeader(apiError.StatusCode())
		return
	}

	metric, apiError := p.Storage.Get(urlData.MetricName, urlData.MetricType)
	if apiError != apierror.NoError {
		w.WriteHeader(apiError.StatusCode())
		return
	}

	if _, err := w.Write([]byte(metric.StringValue())); err != nil {
		log.Printf("Couldn't write response")
	}
}

func (p *Handlers) GetAll(w http.ResponseWriter, _ *http.Request) {
	type KeyValue struct {
		key   string
		value models.Metric
	}

	kv := make([]KeyValue, 0, len(p.Storage.GetAll()))
	for key, value := range p.Storage.GetAll() {
		kv = append(kv, KeyValue{key, value})
	}

	sort.Slice(kv, func(i, j int) bool { return kv[i].key < kv[j].key })
	for _, v := range kv {
		_, err := w.Write([]byte(fmt.Sprintf("Key: %s, value: %s, type: %s \n", v.key, v.value.StringValue(), v.value.ValueType())))
		if err != nil {
			log.Printf("Couldn't write response to GetAll request")
		}
	}
}
