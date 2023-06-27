package main

import (
	"go-metricscol/internal/server"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", server.GetServeMux()))
}
