package configsvc

// Service loads and expands config. Panic occurs in env.Expand (layer 3).
type Service struct{}

// GetDatabaseDSN returns the database DSN after expanding env refs.
// Passes config["dsn"] to Expand; when it's not a string, env.Expand panics.
func (s *Service) GetDatabaseDSN(raw map[string]interface{}) string {
	dsnVal := raw["dsn"]
	return Expand(dsnVal)
}
