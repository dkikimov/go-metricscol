package models

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

//TODO: Точно ли тут float64?

type Gauge float64
type Counter int64

type Metrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	RandomValue   Gauge
	PollCount     Counter
}

// TODO: Насколько Reflect снижает производительность? Как можно сделать иначе?

func (m *Metrics) SendToServer(addr string) {
	iVal := reflect.ValueOf(m).Elem()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		name := iVal.Type().Field(i).Name
		value := fmt.Sprint(f.Interface())

		postUrl := fmt.Sprintf("%s/update/%s/%s/%s", addr, f.Type().Name(), name, value)
		log.Println(postUrl)
		_, err := http.Post(postUrl, "text/plain", nil)
		if err != nil {
			log.Fatalf("Couldn't send metric %s to server", name)
		}
	}
	return
}
