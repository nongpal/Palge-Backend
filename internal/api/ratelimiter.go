package api

import (
	"net/http"

	"golang.org/x/time/rate"
)

func (app *Application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
