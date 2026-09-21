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
	lastSeen   time.Time
}

type RateLimiter struct {
	clientMap sync.Map
	rate      float64
	capacity  float64
	stop      chan struct{}
	stopOnce  sync.Once
}

func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		capacity: capacity,
		stop:     make(chan struct{}),
	}
}

func (rl *RateLimiter) Allow(clientID string) bool {
	now := time.Now()

	actual, _ := rl.clientMap.LoadOrStore(clientID, &ClientInfo{
		tokens:     rl.capacity,
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

	client.lastRefill = now
	client.lastSeen = now

	if client.tokens >= 1 {
		client.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) StartJanitor(interval, expiry time.Duration) {

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				rl.clientMap.Range(func(key, value any) bool {
					client := value.(*ClientInfo)

					client.mu.Lock()
					defer client.mu.Unlock()

					if time.Since(client.lastSeen) > expiry {
						rl.clientMap.Delete(key)
					}
					return true
				})
			case <-rl.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (rl *RateLimiter) StopJanitor() {
	rl.stopOnce.Do(func() {
		close(rl.stop)
	})
}

func (app *Application) middlewareRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.cfg.rl.enabled {
			ip := app.getRealIP(r)
			if !app.rateLimiter.Allow(ip) {
				app.rateLimitExceededResponse(w, r)
				return
			}
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
