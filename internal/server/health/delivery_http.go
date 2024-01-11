package health

import "net/http"

type HTTPHandlers interface {
	Ping(w http.ResponseWriter, r *http.Request)
}
