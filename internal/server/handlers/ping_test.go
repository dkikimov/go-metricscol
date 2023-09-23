package handlers

import (
	"fmt"
	"net/http"
)

func ExampleHandlers_Ping() {
	address := "localhost:8080"

	pingURL := fmt.Sprintf("%s/ping", address)

	http.Get(pingURL)
}
