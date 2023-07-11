package models

import "strconv"

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

func (m MetricType) String() string {
	switch m {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	}
	return ""
}

type Metric struct {
	Name  string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metric) GetStringValue() string {
	switch m.MType {
	case Gauge:
		return strconv.FormatFloat(*m.Value, 'g', -1, 64)
	case Counter:
		return strconv.FormatInt(*m.Delta, 10)
	}

	// TODO: Добавить более строгое ограничение
	return ""
}
