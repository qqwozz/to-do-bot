package config

import (
	"os"
	"testing"
)

func TestLoadDefaultValues(t *testing.T) {
	os.Unsetenv("BOT_TOKEN")
	os.Unsetenv("BACKEND_URL")
	os.Unsetenv("PORT")

	cfg := Load()

	if cfg.BotToken != "" {
		t.Errorf("Expected empty BotToken, got '%s'", cfg.BotToken)
	}
	if cfg.BackendURL != "http://localhost:8081" {
		t.Errorf("Expected 'http://localhost:8081', got '%s'", cfg.BackendURL)
	}
	if cfg.Port != "8080" {
		t.Errorf("Expected '8080', got '%s'", cfg.Port)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("BOT_TOKEN", "test-token-123")
	os.Setenv("BACKEND_URL", "http://backend:9090")
	os.Setenv("PORT", "3000")
	defer os.Unsetenv("BOT_TOKEN")
	defer os.Unsetenv("BACKEND_URL")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.BotToken != "test-token-123" {
		t.Errorf("Expected 'test-token-123', got '%s'", cfg.BotToken)
	}
	if cfg.BackendURL != "http://backend:9090" {
		t.Errorf("Expected 'http://backend:9090', got '%s'", cfg.BackendURL)
	}
	if cfg.Port != "3000" {
		t.Errorf("Expected '3000', got '%s'", cfg.Port)
	}
}

func TestLoadPartialEnv(t *testing.T) {
	os.Setenv("BOT_TOKEN", "my-token")
	os.Unsetenv("BACKEND_URL")
	os.Unsetenv("PORT")
	defer os.Unsetenv("BOT_TOKEN")

	cfg := Load()

	if cfg.BotToken != "my-token" {
		t.Errorf("Expected 'my-token', got '%s'", cfg.BotToken)
	}
	if cfg.BackendURL != "http://localhost:8081" {
		t.Errorf("Expected default BackendURL")
	}
}
