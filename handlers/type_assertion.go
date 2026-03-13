package handlers

import (
	"net/http"
)

// ConfigLoader reads configuration from a map and extracts typed sections.
// The bug: GetDatabaseConfig performs an unsafe type assertion — config["database"]
// is a string, not map[string]interface{}, so the assertion panics.
type ConfigLoader struct {
	raw map[string]interface{}
}

// NewConfigLoader returns a loader with the given raw config.
func NewConfigLoader(raw map[string]interface{}) *ConfigLoader {
	return &ConfigLoader{raw: raw}
}

// GetDatabaseConfig extracts database host, port, and name from the config.
// It assumes config["database"] is a map[string]interface{}; when it is actually
// a string (e.g. a DSN or env var name), the type assertion panics.
func (c *ConfigLoader) GetDatabaseConfig() (host string, port int, name string) {
	dbSection := c.raw["database"].(map[string]interface{})
	host = dbSection["host"].(string)
	port = dbSection["port"].(int)
	name = dbSection["name"].(string)
	return host, port, name
}

// TypeAssertion handles GET /error/type-assertion.
// It uses a config where "database" is a string, so GetDatabaseConfig panics.
func TypeAssertion(loader *ConfigLoader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In real code this might come from YAML/JSON; here we set database to a string.
		config := map[string]interface{}{
			"database": "postgres://localhost/mydb", // string, not a map
			"app":      map[string]interface{}{"name": "error-simulator"},
		}
		l := NewConfigLoader(config)
		_, _, _ = l.GetDatabaseConfig()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
