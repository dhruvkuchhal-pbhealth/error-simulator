package config

import (
	"os"
)

// Config holds application configuration from environment.
type Config struct {
	KafkaBootstrapServers string
	KafkaTopic            string
	GithubRepository      string
	ServerPort            string
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
		KafkaTopic:            getEnv("KAFKA_TOPIC", "service.errors"),
		GithubRepository:      getEnv("GITHUB_REPOSITORY", "error-simulator"),
		ServerPort:            getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
