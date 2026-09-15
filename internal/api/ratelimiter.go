package api

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type ClientInfo struct {
	mu       sync.Mutex
	tokens   int
	lastSeen time.Time
}

type RateLimiter struct {
	clientMap sync.Map
	rate      int
	window    time.Duration
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:   rate,
		window: window,
	}
}

func (rl *RateLimiter) Allow(clientID string) bool {
	now := time.Now()

	actual, _ := rl.clientMap.LoadOrStore(clientID, &ClientInfo{
		tokens:   rl.rate,
		lastSeen: now,
	})

	client := actual.(*ClientInfo)

	client.mu.Lock()
	defer client.mu.Unlock()

	if now.Sub(client.lastSeen) >= rl.window {
		client.tokens = rl.rate
		client.lastSeen = now
	}

	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false

}
func (app *Application) middlewareRateLimit(next http.Handler) http.Handler {
	limiter := NewRateLimiter(10, 10*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := app.getRealIP(r)
		if !limiter.Allow(ip) {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) getRealIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
