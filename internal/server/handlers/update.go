package handlers

import (
	"log"
	"net/http"
	"strings"
)

func UpdateGauge(w http.ResponseWriter, r *http.Request) {
	data := strings.Split(r.URL.Path, "/")[3:]
	log.Println(data)
}

func UpdateCounter(w http.ResponseWriter, r *http.Request) {
	data := strings.Split(r.URL.Path, "/")[3:]
	log.Println(data)
}
