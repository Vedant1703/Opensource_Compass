package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort     string
	DatabaseURL    string
	GitHubSvcURL   string
	AISvcURL       string
	JWTSecret 	   string
}

func Load() *Config {
	cfg := &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		GitHubSvcURL: getEnv("GITHUB_SERVICE_URL", "http://localhost:8081"),
		AISvcURL:     getEnv("AI_SERVICE_URL", "http://localhost:8082"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret"),
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultVal == "" {
		log.Printf("WARNING: %s not set\n", key)
	}

	return defaultVal
}
