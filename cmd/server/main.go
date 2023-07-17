package main

import (
	"go-metricscol/internal/server"
	"log"
)

func main() {
	srv := server.New("127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}
