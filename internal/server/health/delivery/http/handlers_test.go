package http

import (
	"fmt"
	"net/http"
)

func ExampleHealthHandlers_Ping() {
	address := "localhost:8080"

	pingURL := fmt.Sprintf("%s/ping", address)

	response, err := http.Get(pingURL)
	if err != nil {
		// Handle error
		return
	}
	response.Body.Close()
}
