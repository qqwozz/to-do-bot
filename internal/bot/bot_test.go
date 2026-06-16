package bot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"todo-bot/internal/config"
)

func TestNewWithInvalidToken(t *testing.T) {
	cfg := &config.Config{
		BotToken:   "invalid:token",
		BackendURL: "http://localhost:8081",
	}

	_, err := New(cfg)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestNewClientInitialized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		BotToken:   "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BackendURL: server.URL,
	}

	bot, err := New(cfg)
	if err != nil {
		t.Skipf("Skipping bot creation test (network): %v", err)
	}

	if bot.client == nil {
		t.Error("Client should be initialized")
	}
	if bot.config != cfg {
		t.Error("Config should be stored")
	}
}

func TestNewClientWithCustomBackend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	if client.baseURL != server.URL {
		t.Errorf("Expected baseURL %s, got %s", server.URL, client.baseURL)
	}
	if client.httpClient == nil {
		t.Error("HTTP client should be initialized")
	}
}

func TestCreatePlanRequestSerialization(t *testing.T) {
	plan := PlanRequest{
		Title:       "Тест",
		Description: "Описание",
		Date:        "2024-12-25",
		Time:        "14:00",
		IsAllDay:    false,
	}

	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal PlanRequest: %v", err)
	}

	var decoded PlanRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal PlanRequest: %v", err)
	}

	if decoded.Title != plan.Title {
		t.Errorf("Title mismatch: %s != %s", decoded.Title, plan.Title)
	}
	if decoded.Description != plan.Description {
		t.Errorf("Description mismatch: %s != %s", decoded.Description, plan.Description)
	}
	if decoded.Date != plan.Date {
		t.Errorf("Date mismatch: %s != %s", decoded.Date, plan.Date)
	}
	if decoded.Time != plan.Time {
		t.Errorf("Time mismatch: %s != %s", decoded.Time, plan.Time)
	}
	if decoded.IsAllDay != plan.IsAllDay {
		t.Errorf("IsAllDay mismatch: %v != %v", decoded.IsAllDay, plan.IsAllDay)
	}
}

func TestPlanSerialization(t *testing.T) {
	plan := Plan{
		ID:          1,
		Title:       "Тест",
		Description: "Описание",
		Date:        "2024-12-25",
		Time:        "14:00",
		IsAllDay:    false,
		CreatedAt:   "2024-12-24T10:00:00Z",
	}

	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal Plan: %v", err)
	}

	var decoded Plan
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Plan: %v", err)
	}

	if decoded.ID != plan.ID {
		t.Errorf("ID mismatch: %d != %d", decoded.ID, plan.ID)
	}
	if decoded.Title != plan.Title {
		t.Errorf("Title mismatch: %s != %s", decoded.Title, plan.Title)
	}
	if decoded.CreatedAt != plan.CreatedAt {
		t.Errorf("CreatedAt mismatch: %s != %s", decoded.CreatedAt, plan.CreatedAt)
	}
}

func TestPlanRequestJSONTags(t *testing.T) {
	plan := PlanRequest{
		Title:       "Заголовок",
		Description: "Описание дела",
		Date:        "2025-01-15",
		Time:        "09:30",
		IsAllDay:    true,
	}

	data, _ := json.Marshal(plan)
	jsonStr := string(data)

	expectedFields := []string{`"title"`, `"description"`, `"date"`, `"time"`, `"is_all_day"`}
	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON should contain field %s, got: %s", field, jsonStr)
		}
	}
}
