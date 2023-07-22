package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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

func (m *Metric) UnmarshalJSON(bytes []byte) error {
	type MetricAlias Metric

	var tempMetric MetricAlias
	err := json.Unmarshal(bytes, &tempMetric)
	if err != nil {
		return err
	}

	if len(tempMetric.Name) == 0 || len(tempMetric.MType) == 0 {
		return fmt.Errorf("empty name or type")
	}

	if (tempMetric.MType == Gauge && tempMetric.Value == nil) || (tempMetric.MType == Counter && tempMetric.Delta == nil) {
		return fmt.Errorf("empty value")
	}

	*m = (Metric)(tempMetric)
	return nil
}

func (m *Metric) GetStringValue() string {
	switch m.MType {
	case Gauge:
		return strconv.FormatFloat(*m.Value, 'g', -1, 64)
	case Counter:
		return strconv.FormatInt(*m.Delta, 10)
	}

	// TODO: Добавить более строгое ограничение
	return ""
}
