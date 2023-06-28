package main

import (
	"go-metricscol/internal/server"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", server.GetServeMux()))
}
