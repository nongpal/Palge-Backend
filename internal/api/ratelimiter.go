package api

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type ClientInfo struct {
	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
}

type RateLimiter struct {
	clientMap sync.Map
	rate      float64
	capacity  float64
}

func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		capacity: capacity,
	}
}

func (rl *RateLimiter) Allow(clientID string) bool {
	now := time.Now()

	actual, _ := rl.clientMap.LoadOrStore(clientID, &ClientInfo{
		tokens:     rl.rate,
		lastRefill: now,
	})

	client := actual.(*ClientInfo)

	client.mu.Lock()
	defer client.mu.Unlock()

	elapsed := now.Sub(client.lastRefill).Seconds()
	client.tokens += elapsed * rl.rate

	if client.tokens > rl.capacity {
		client.tokens = rl.capacity
	}

	if client.tokens >= 1 {
		client.tokens--
		return true
	}

	return false

}
func (app *Application) middlewareRateLimit(next http.Handler) http.Handler {
	limiter := NewRateLimiter(2, 10)

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
