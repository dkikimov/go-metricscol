package health

import "net/http"

type HttpHandlers interface {
	Ping(w http.ResponseWriter, r *http.Request)
}
