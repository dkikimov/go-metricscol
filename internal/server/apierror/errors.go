//TODO: Как лучше назвать пакет?

package apierror

import (
	"net/http"
)

type APIError struct {
	StatusCode int
	Message    string
}

// Write custom api error
func (apiError APIError) Error() string {
	return apiError.Message
}

func WriteHeader(w http.ResponseWriter, err error) {
	if apiError, ok := err.(APIError); ok {
		w.WriteHeader(apiError.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//TODO: Стоит ли писать тесты для подобных файлов?

var (
	NotEnoughArguments = APIError{
		StatusCode: http.StatusNotImplemented,
		Message:    "not enough arguments",
	}
	UnknownMetricType = APIError{
		StatusCode: http.StatusNotImplemented,
		Message:    "unknown metric type",
	}
	EmptyArguments = APIError{
		StatusCode: http.StatusNotFound,
		Message:    "empty arguments",
	}

	InvalidValue = APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid value",
	}

	NumberParse = APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "couldn't parse number",
	}

	NotFound = APIError{
		StatusCode: http.StatusNotFound,
		Message:    "not found",
	}
)
