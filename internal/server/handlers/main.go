package handlers

import (
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
}
