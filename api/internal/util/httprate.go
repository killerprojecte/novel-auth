package util

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func RateLimiter(limit int) func(next http.Handler) http.Handler {
	return httprate.Limit(limit, time.Hour,
		httprate.WithKeyFuncs(httprate.KeyByRealIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "操作过于频繁，请稍后再试", http.StatusTooManyRequests)
		}),
	)
}
