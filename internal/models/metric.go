package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	Hash  string     `json:"hash,omitempty"`  // значение хеш-функции
}

func (m Metric) StringValue() string {
	switch m.MType {
	case Gauge:
		return strconv.FormatFloat(*m.Value, 'g', -1, 64)
	case Counter:
		return strconv.FormatInt(*m.Delta, 10)
	}

	// TODO: Добавить более строгое ограничение
	return ""
}

func (m Metric) HashValue(id string) string {
	if len(id) == 0 {
		return ""
	}

	h := hmac.New(sha256.New, []byte(id))
	var str string
	switch m.MType {
	case Counter:
		str = fmt.Sprintf("%s:counter:%d", m.Name, *m.Delta)
		break
	case Gauge:
		str = fmt.Sprintf("%s:gauge:%f", m.Name, *m.Value)
		break
	default:
		return ""
	}

	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
