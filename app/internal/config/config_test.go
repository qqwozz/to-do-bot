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
	if cfg.Port != "8080" {
		t.Errorf("Expected default Port '8080', got '%s'", cfg.Port)
	}
}

func TestLoadOnlyBackendURL(t *testing.T) {
	os.Unsetenv("BOT_TOKEN")
	os.Setenv("BACKEND_URL", "http://custom:8080")
	os.Unsetenv("PORT")
	defer os.Unsetenv("BACKEND_URL")

	cfg := Load()

	if cfg.BotToken != "" {
		t.Errorf("Expected empty BotToken, got '%s'", cfg.BotToken)
	}
	if cfg.BackendURL != "http://custom:8080" {
		t.Errorf("Expected 'http://custom:8080', got '%s'", cfg.BackendURL)
	}
	if cfg.Port != "8080" {
		t.Errorf("Expected default Port '8080', got '%s'", cfg.Port)
	}
}

func TestLoadOnlyPort(t *testing.T) {
	os.Unsetenv("BOT_TOKEN")
	os.Unsetenv("BACKEND_URL")
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.BotToken != "" {
		t.Errorf("Expected empty BotToken, got '%s'", cfg.BotToken)
	}
	if cfg.BackendURL != "http://localhost:8081" {
		t.Errorf("Expected default BackendURL, got '%s'", cfg.BackendURL)
	}
	if cfg.Port != "9090" {
		t.Errorf("Expected '9090', got '%s'", cfg.Port)
	}
}

func TestLoadEmptyValues(t *testing.T) {
	os.Unsetenv("BOT_TOKEN")
	os.Unsetenv("BACKEND_URL")
	os.Unsetenv("PORT")

	cfg := Load()

	if cfg.BotToken != "" {
		t.Errorf("Expected empty BotToken, got '%s'", cfg.BotToken)
	}
	if cfg.BackendURL != "http://localhost:8081" {
		t.Errorf("Expected default BackendURL, got '%s'", cfg.BackendURL)
	}
	if cfg.Port != "8080" {
		t.Errorf("Expected default Port, got '%s'", cfg.Port)
	}
}

func TestLoadSpecialCharacters(t *testing.T) {
	os.Setenv("BOT_TOKEN", "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
	os.Setenv("BACKEND_URL", "http://user:pass@backend:9090/path")
	os.Setenv("PORT", "443")
	defer os.Unsetenv("BOT_TOKEN")
	defer os.Unsetenv("BACKEND_URL")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.BotToken != "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11" {
		t.Errorf("Expected special token, got '%s'", cfg.BotToken)
	}
	if cfg.BackendURL != "http://user:pass@backend:9090/path" {
		t.Errorf("Expected URL with auth, got '%s'", cfg.BackendURL)
	}
	if cfg.Port != "443" {
		t.Errorf("Expected '443', got '%s'", cfg.Port)
	}
}

func TestLoadReturnPointer(t *testing.T) {
	cfg := Load()
	if cfg == nil {
		t.Error("Load should return non-nil Config")
	}
}

func TestLoadConfigFields(t *testing.T) {
	os.Setenv("BOT_TOKEN", "token")
	os.Setenv("BACKEND_URL", "http://backend:8081")
	os.Setenv("PORT", "8080")
	defer os.Unsetenv("BOT_TOKEN")
	defer os.Unsetenv("BACKEND_URL")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.BotToken == "" {
		t.Error("BotToken should not be empty")
	}
	if cfg.BackendURL == "" {
		t.Error("BackendURL should not be empty")
	}
	if cfg.Port == "" {
		t.Error("Port should not be empty")
	}
}

func TestGetEnvExisting(t *testing.T) {
	os.Setenv("TEST_GETENV_KEY", "test-value")
	defer os.Unsetenv("TEST_GETENV_KEY")

	result := getEnv("TEST_GETENV_KEY", "default")
	if result != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", result)
	}
}

func TestGetEnvNotExisting(t *testing.T) {
	os.Unsetenv("TEST_GETENV_MISSING")

	result := getEnv("TEST_GETENV_MISSING", "default-value")
	if result != "default-value" {
		t.Errorf("Expected 'default-value', got '%s'", result)
	}
}

func TestGetEnvEmptyValue(t *testing.T) {
	os.Setenv("TEST_GETENV_EMPTY", "")
	defer os.Unsetenv("TEST_GETENV_EMPTY")

	result := getEnv("TEST_GETENV_EMPTY", "default")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}
