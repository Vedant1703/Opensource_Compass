package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// RequestID middleware generates or reads an X-Request-ID header and stores it in context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Store in context for logging
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Echo the request ID back in the response header
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from a context.
func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}

// Logf is a simple helper to prefix every log line with the request ID.
func Logf(ctx context.Context, format string, args ...interface{}) {
	reqID := GetRequestID(ctx)
	if reqID != "" {
		log.Printf("[req:%s] "+format, append([]interface{}{reqID[:8]}, args...)...)
	} else {
		log.Printf(format, args...)
	}
}
