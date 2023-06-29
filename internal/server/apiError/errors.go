//TODO: Как лучше назвать пакет?

package apiError

import "net/http"

type APIError int

func (err APIError) StatusCode() int {
	return int(err)
}

const (
	NotEnoughArguments = http.StatusNotFound
	UnknownMetricType  = http.StatusNotImplemented
	EmptyArguments     = http.StatusNotFound
	NoError            = http.StatusOK
	TypeMismatch       = http.StatusBadRequest
	NumberParse        = http.StatusBadRequest
)
