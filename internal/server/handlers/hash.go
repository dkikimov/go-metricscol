package handlers

import (
	"go-metricscol/internal/models"
)

func (p *Handlers) addHash(metric *models.Metric) {
	metric.Hash = metric.HashValue(p.config.HashKey)
}

func (p *Handlers) addHashToSlice(all []models.Metric) {
	for idx, value := range all {
		all[idx].Hash = value.HashValue(p.config.HashKey)
	}
}
