package config

import "os"

type Config struct {
	Port      string
	GeminiKey string
	Model     string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("AI_SERVICE_PORT", "8081"),
		GeminiKey: os.Getenv("GEMINI_API_KEY"),
		Model: getEnv("LLM_MODEL", "models/gemini-2.5-flash"),

	}
}


func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
