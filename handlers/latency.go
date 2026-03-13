package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/your-org/error-simulator/config"
	"github.com/your-org/error-simulator/logger"
)

const (
	maxLatencyMs = 300_000 // 5 minutes cap
)

// Latency handles GET /error/latency. Delays response by a configurable duration (env LATENCY_MS or ?ms=).
// No panic — use for testing timeouts, slow clients, or latency monitoring.
func Latency(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ms := cfg.LatencyMs
		if q := r.URL.Query().Get("ms"); q != "" {
			if n, err := strconv.Atoi(q); err == nil && n >= 0 {
				if n > maxLatencyMs {
					n = maxLatencyMs
				}
				ms = n
			}
		}
		d := time.Duration(ms) * time.Millisecond
		logger.Log.Info().Int("ms", ms).Str("path", r.URL.Path).Msg("latency: sleeping before response")
		time.Sleep(d)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
