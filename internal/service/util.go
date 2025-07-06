package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

func NotFound(message string) *HttpError {
	return &HttpError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func BadRequest(message string) *HttpError {
	return &HttpError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func Unauthorized(message string) *HttpError {
	return &HttpError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func InternalServerError(message string) *HttpError {
	return &HttpError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func ToHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			httpErr := &HttpError{}
			if errors.As(err, &httpErr) {
				http.Error(w, httpErr.Message, httpErr.StatusCode)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func body[T any](r *http.Request) (T, error) {
	var zero T

	contentType := r.Header.Get("Content-Type")
	if contentType != "" && contentType != "application/json" {
		return zero, &HttpError{
			StatusCode: http.StatusUnsupportedMediaType,
			Message:    "expected content-type application/json",
		}
	}

	// 限制读取的最大字节数为1MB
	const maxBytesDefault = 1 << 20
	limitedReader := io.LimitReader(r.Body, maxBytesDefault)
	defer r.Body.Close()

	// 解码JSON
	var result T
	if err := json.NewDecoder(limitedReader).Decode(&result); err != nil {
		return zero, &HttpError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
	}

	return result, nil
}

func Respond[T any](w http.ResponseWriter, statusCode int, response T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return InternalServerError("failed to encode response")
	}
	return nil
}
