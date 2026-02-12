package auth

import (
	"crypto/subtle"
	"net/http"
)

// APIKeyAuth returns middleware that checks for a valid API key
// in the Authorization header (Bearer <key>) or X-API-Key header.
// If validKey is empty, auth is disabled (development mode).
func APIKeyAuth(validKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth if no key is configured (local dev)
			if validKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Allow OPTIONS for CORS preflight
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			key := extractKey(r)
			if key == "" {
				http.Error(w, `{"error":"missing API key"}`, http.StatusUnauthorized)
				return
			}

			if subtle.ConstantTimeCompare([]byte(key), []byte(validKey)) != 1 {
				http.Error(w, `{"error":"invalid API key"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractKey tries Authorization: Bearer <key>, then X-API-Key header.
func extractKey(r *http.Request) string {
	// Check Authorization header
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}

	// Check X-API-Key header
	return r.Header.Get("X-API-Key")
}
