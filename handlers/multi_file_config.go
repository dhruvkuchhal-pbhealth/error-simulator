package handlers

import (
	"net/http"

	"github.com/your-org/error-simulator/configsvc"
)

// MultiFileConfig handles GET /error/multi-file/config.
// Call chain: handler → configsvc.GetDatabaseDSN → env.Expand (panic in env).
func MultiFileConfig(svc *configsvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := map[string]interface{}{
			"dsn": map[string]string{"host": "localhost"}, // not a string → env.Expand panics
		}
		_ = svc.GetDatabaseDSN(raw)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
