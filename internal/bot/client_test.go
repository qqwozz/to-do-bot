package bot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8081")
	if client.baseURL != "http://localhost:8081" {
		t.Errorf("Expected baseURL 'http://localhost:8081', got '%s'", client.baseURL)
	}
}

func TestCreatePlanSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/plans" {
			t.Errorf("Expected /plans, got %s", r.URL.Path)
		}

		var plan PlanRequest
		json.NewDecoder(r.Body).Decode(&plan)

		if plan.Title == "" {
			t.Error("Title should not be empty")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      1,
			"message": "План создан",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	plan := PlanRequest{
		Title:       "Тест",
		Description: "Описание",
		Date:        "2024-12-25",
		Time:        "14:00",
		IsAllDay:    false,
	}

	err := client.CreatePlan(plan)
	if err != nil {
		t.Errorf("CreatePlan returned error: %v", err)
	}
}

func TestCreatePlanError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	plan := PlanRequest{Title: "Тест", Date: "2024-12-25"}

	err := client.CreatePlan(plan)
	if err == nil {
		t.Error("Expected error from CreatePlan")
	}
}

func TestGetPlansByDate(t *testing.T) {
	plans := []Plan{
		{ID: 1, Title: "План 1", Date: "2024-12-25"},
		{ID: 2, Title: "План 2", Date: "2024-12-25"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("date") != "2024-12-25" {
			t.Error("Expected date parameter")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetPlansByDate("2024-12-25")

	if err != nil {
		t.Errorf("GetPlansByDate returned error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 plans, got %d", len(result))
	}
	if result[0].Title != "План 1" {
		t.Errorf("Expected 'План 1', got '%s'", result[0].Title)
	}
}

func TestGetPlansByDateRange(t *testing.T) {
	plans := []Plan{
		{ID: 1, Title: "План 1", Date: "2024-12-23"},
		{ID: 2, Title: "План 2", Date: "2024-12-25"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		if start != "2024-12-23" || end != "2024-12-29" {
			t.Errorf("Expected start=2024-12-23, end=2024-12-29, got start=%s, end=%s", start, end)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetPlansByDateRange("2024-12-23", "2024-12-29")

	if err != nil {
		t.Errorf("GetPlansByDateRange returned error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 plans, got %d", len(result))
	}
}

func TestGetPlansByDateError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPlansByDate("2024-12-25")

	if err == nil {
		t.Error("Expected error from GetPlansByDate")
	}
}

func TestGetJSONInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var result []Plan
	err := client.getJSON("/plans", &result)

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
