package config

import "os"

type Config struct {
	BotToken   string
	BackendURL string
	Port       string
}

func Load() *Config {
	return &Config{
		BotToken:   getEnv("BOT_TOKEN", ""),
		BackendURL: getEnv("BACKEND_URL", "http://localhost:8081"),
		Port:       getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
