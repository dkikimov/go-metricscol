package server

import (
	"net/http"
)

func GetServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
