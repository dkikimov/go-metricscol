package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strconv"
)

// MetricType is type describing type of metric.
type MetricType string

// Declaration of the metric types used.
const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

// String returns string representation.
func (m MetricType) String() string {
	switch m {
	case Gauge:
		return "gauge"
	case Counter:
		return "counter"
	}
	return ""
}

// Scan overrides logic of scanning MetricType in database.
func (m *MetricType) Scan(src any) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("can't convert %T to string", src)
	}

	switch str {
	case Gauge.String():
		*m = Gauge
	case Counter.String():
		*m = Counter
	default:
		return fmt.Errorf("unknown metric type %s", src)
	}
	return nil
}

// Value overrides logic of storing MetricType in database.
func (m *MetricType) Value() (driver.Value, error) {
	return m.String(), nil
}

// Metric is a description of metric entity.
type Metric struct {
	Name  string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string     `json:"hash,omitempty"`  // значение хеш-функции
}

// StringValue returns metric value in string.
func (m *Metric) StringValue() string {
	switch m.MType {
	case Gauge:
		return strconv.FormatFloat(*m.Value, 'g', -1, 64)
	case Counter:
		return strconv.FormatInt(*m.Delta, 10)
	}

	// TODO: Добавить более строгое ограничение
	return ""
}

// HashValue returns hash of metric based on name, type and value.
func (m *Metric) HashValue(id string) string {
	if len(id) == 0 {
		return ""
	}

	h := hmac.New(sha256.New, []byte(id))
	var str string
	switch m.MType {
	case Counter:
		str = fmt.Sprintf("%s:counter:%d", m.Name, *m.Delta)
	case Gauge:
		str = fmt.Sprintf("%s:gauge:%f", m.Name, *m.Value)
	default:
		return ""
	}

	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
