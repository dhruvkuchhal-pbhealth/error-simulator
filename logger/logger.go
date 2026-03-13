package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Log is the global logger. Human-readable console output so "what fucked up"
// is obvious when triaging errors and raising fix PRs.
var Log zerolog.Logger

func init() {
	outputs := []io.Writer{zerolog.ConsoleWriter{Out: os.Stdout}}
	if w := initLogstashOutput(); w.host != "" {
		outputs = append(outputs, w)
	}

	Log = zerolog.New(io.MultiWriter(outputs...)).
		With().
		Timestamp().
		Str("service", "error-simulator").
		Str("layer", "simulator").
		Logger()
}
