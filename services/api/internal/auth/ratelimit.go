package auth

import (
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	count       int
	windowStart time.Time
}

// RateLimiter returns middleware that limits requests per IP per minute.
// If maxPerMinute <= 0, rate limiting is disabled.
func RateLimiter(maxPerMinute int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	// Clean up old entries every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			now := time.Now()
			for ip, v := range visitors {
				if now.Sub(v.windowStart) > 2*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if maxPerMinute <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			ip := r.RemoteAddr
			// Strip port if present
			if idx := len(ip) - 1; idx > 0 {
				for i := idx; i >= 0; i-- {
					if ip[i] == ':' {
						ip = ip[:i]
						break
					}
				}
			}

			mu.Lock()
			now := time.Now()
			v, exists := visitors[ip]
			if !exists || now.Sub(v.windowStart) > time.Minute {
				visitors[ip] = &visitor{count: 1, windowStart: now}
				mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			v.count++
			if v.count > maxPerMinute {
				mu.Unlock()
				w.Header().Set("Retry-After", "60")
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
