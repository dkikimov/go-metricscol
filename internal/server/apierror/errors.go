//TODO: Как лучше назвать пакет?

package apierror

import "net/http"

type APIError int

//TODO: Стоит ли писать тесты для подобных файлов?

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
	NotFound           = http.StatusNotFound
)
