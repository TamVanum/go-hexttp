package hexttp

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"
const requestStart contextKey = "requestStart"

func MetadataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		ctx = context.WithValue(ctx, requestStart, start)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}
