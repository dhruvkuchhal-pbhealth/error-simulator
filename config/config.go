package config

import (
	"os"
	"strconv"
)

// Config holds application configuration from environment.
type Config struct {
	KafkaBootstrapServers string
	KafkaTopic            string
	GithubRepository      string
	ServerPort            string
	// LatencyMs is the default delay (ms) for the latency endpoint. Overridable via ?ms= query.
	LatencyMs int
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaTopic:            getEnv("KAFKA_TOPIC", "service.errors"),
		GithubRepository:      getEnv("GITHUB_REPOSITORY", "error-simulator"),
		ServerPort:            getEnv("SERVER_PORT", "8080"),
		LatencyMs:             getEnvInt("LATENCY_MS", 3000),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			return n
		}
	}
	return defaultVal
}
