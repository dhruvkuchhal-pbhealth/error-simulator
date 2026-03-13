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
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "10.0.10.135:9092"),
		KafkaTopic:            getEnv("KAFKA_TOPIC", "app-error-logs"),
		GithubRepository:      getEnv("GITHUB_REPOSITORY", "dhruvkuchhal-pbhealth/error-simulator"),
		ServerPort:            getEnv("SERVER_PORT", "8092"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
