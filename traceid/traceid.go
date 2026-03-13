package traceid

import "context"

type key struct{}

// WithContext returns a context that carries the given trace ID.
func WithContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, key{}, traceID)
}

// FromContext returns the trace ID from ctx if set.
func FromContext(ctx context.Context) string {
	v := ctx.Value(key{})
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
