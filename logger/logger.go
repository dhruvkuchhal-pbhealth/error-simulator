package logger

import (
	"os"

	"github.com/rs/zerolog"
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
