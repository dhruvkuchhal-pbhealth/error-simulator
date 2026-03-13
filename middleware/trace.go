package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/your-org/error-simulator/traceid"
)

const (
	// HeaderTraceID is the request/response header for trace ID (e.g. from gateways or clients).
	HeaderTraceID = "X-Trace-ID"
	// HeaderRequestID is an alternative header some clients send.
	HeaderRequestID = "X-Request-ID"
)

// TraceID returns a middleware that injects a trace ID into the request context
// and sets it on the response so callers can correlate logs and events across
// functions and repos. It reads X-Trace-ID or X-Request-ID from the request;
// if missing, it generates a new ID.
func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(HeaderTraceID)
		if traceID == "" {
			traceID = r.Header.Get(HeaderRequestID)
		}
		if traceID == "" {
			traceID = generateTraceID()
		}
		traceID = strings.TrimSpace(traceID)
		if traceID == "" {
			traceID = generateTraceID()
		}
		ctx := traceid.WithContext(r.Context(), traceID)
		w.Header().Set(HeaderTraceID, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TraceIDFromContext returns the trace ID from ctx if set by TraceID middleware.
func TraceIDFromContext(ctx context.Context) string {
	return traceid.FromContext(ctx)
}

func generateTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback-" + hex.EncodeToString(b[:8]) // b may be zero; still unique enough
	}
	return hex.EncodeToString(b)
}
