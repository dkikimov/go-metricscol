package models

import (
	"strconv"
)

type MetricType int

const (
	GaugeType   MetricType = iota // float64
	CounterType                   //int64
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

type Metric interface {
	GetName() string
	GetType() MetricType
	GetStringValue() string
}

type Gauge struct {
	Name  string
	Value float64
}

func (g Gauge) GetName() string {
	return g.Name
}

func (g Gauge) GetType() MetricType {
	return GaugeType
}

func (g Gauge) GetStringValue() string {
	return strconv.FormatFloat(g.Value, 'g', -1, 64)
}

type Counter struct {
	Name  string
	Value int64
}

func (c Counter) GetName() string {
	return c.Name
}

func (c Counter) GetType() MetricType {
	return CounterType
}

func (c Counter) GetStringValue() string {
	return strconv.FormatInt(c.Value, 10)
}
