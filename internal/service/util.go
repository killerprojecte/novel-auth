package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type httpError struct {
	StatusCode int
	Message    string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

func notFound(message string) *httpError {
	return &httpError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func badRequest(message string) *httpError {
	return &httpError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func unauthorized(message string) *httpError {
	return &httpError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func internalServerError(message string) *httpError {
	return &httpError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func toHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			httpErr := &httpError{}
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
		return zero, &httpError{
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
		return zero, badRequest("invalid JSON format")
	}

	// 验证JSON
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(result); err != nil {
		errors := err.(validator.ValidationErrors)
		return zero, badRequest(fmt.Sprintf("invalid request: %s", errors))
	}

	return result, nil
}

func respond[T any](w http.ResponseWriter, statusCode int, response T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return internalServerError("failed to encode response")
	}
	return nil
}
