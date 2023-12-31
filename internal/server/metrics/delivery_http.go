package metrics

import (
	"net/http"
)

type HttpHandlers interface {
	Find(w http.ResponseWriter, r *http.Request)
	FindJSON(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	UpdateJSON(w http.ResponseWriter, r *http.Request)
	Updates(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}
