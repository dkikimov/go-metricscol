package apierror

import (
	"errors"
	"net/http"
)

type APIError struct {
	StatusCode int
	Message    string
}

func (apiError APIError) Error() string {
	return apiError.Message
}

func WriteHTTP(w http.ResponseWriter, err error) {
	var apiError APIError
	if errors.As(err, &apiError) {
		w.WriteHeader(apiError.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

var (
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
