package util

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

func NewHttpError(statusCode int, message string) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func NotFound(message string) *HttpError {
	return NewHttpError(http.StatusNotFound, message)
}

func BadRequest(message string) *HttpError {
	return NewHttpError(http.StatusBadRequest, message)
}

func Unauthorized(message string) *HttpError {
	return NewHttpError(http.StatusUnauthorized, message)
}

func InternalServerError(message string) *HttpError {
	return NewHttpError(http.StatusInternalServerError, message)
}

func E(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			httpErr := &HttpError{}
			if errors.As(err, &httpErr) {
				RespondError(w, httpErr.StatusCode, httpErr.Message)
			} else {
				RespondError(w, http.StatusInternalServerError, err.Error())
			}
		}
	}
}
