package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 修复http.Error的额外换行符问题
func RespondError(w http.ResponseWriter, code int, message string) {
	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

func RespondText(w http.ResponseWriter, message string) error {
	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, message)
	return nil
}

func RespondJson[T any](w http.ResponseWriter, response T) error {
	h := w.Header()
	h.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return InternalServerError("failed to encode response")
	}
	return nil
}
