package middleware

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

// RequestID generates a unique request ID and stores it in the
// X-Request-Id response header. If the incoming request already
// carries the header, its value is reused.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-Id")
		if id == "" {
			id = generateID()
		}
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r)
	})
}

// generateID produces a random 16-byte hex string.
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
