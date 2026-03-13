package configsvc

// Expand expands a config value (e.g. env var reference) into a string.
// BUG: Assumes v is always a string; when v is a map (e.g. nested config), type assertion panics (multi-file stack).
func Expand(v interface{}) string {
	return v.(string)
}
