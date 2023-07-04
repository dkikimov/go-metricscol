package models

import "strconv"

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType            = "counter"
)

func (m MetricType) String() string {
	switch m {
	case GaugeType:
		return "gauge"
	case CounterType:
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
	case GaugeType:
		return strconv.FormatFloat(*m.Value, 'g', -1, 64)
	case CounterType:
		return strconv.FormatInt(*m.Delta, 10)
	}

	// TODO: Добавить более строгое ограничение
	return ""
}

//type Metric interface {
//	GetName() string
//	GetType() MetricType
//	GetStringValue() string
//}
//
//type Gauge struct {
//	Name  string
//	Value float64
//}
//
//func (g Gauge) GetName() string {
//	return g.Name
//}
//
//func (g Gauge) GetType() MetricType {
//	return GaugeType
//}
//
//func (g Gauge) GetStringValue() string {
//	return strconv.FormatFloat(g.Value, 'g', -1, 64)
//}
//
//type Counter struct {
//	Name  string
//	Value int64
//}
//
//func (c Counter) GetName() string {
//	return c.Name
//}
//
//func (c Counter) GetType() MetricType {
//	return CounterType
//}
//
//func (c Counter) GetStringValue() string {
//	return strconv.FormatInt(c.Value, 10)
//}
