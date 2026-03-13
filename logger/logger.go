package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/your-org/error-simulator/traceid"
)

// Log is the global logger. Human-readable console output so "what fucked up"
// is obvious when triaging errors and raising fix PRs.
var Log zerolog.Logger

func init() {
	Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		Str("service", "error-simulator").
		Logger()
}

// WithTrace returns a logger that includes trace_id from ctx when set by middleware.TraceID.
// Use for request-scoped logs so traces can be followed across functions and repos.
func WithTrace(ctx context.Context) *zerolog.Logger {
	traceID := traceid.FromContext(ctx)
	if traceID == "" {
		return &Log
	}
	l := Log.With().Str("trace_id", traceID).Logger()
	return &l
}
