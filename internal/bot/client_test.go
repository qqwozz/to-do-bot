package bot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:8081")
	if client.baseURL != "http://localhost:8081" {
		t.Errorf("Expected baseURL 'http://localhost:8081', got '%s'", client.baseURL)
	}
	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
}

func TestNewClientTimeout(t *testing.T) {
	client := NewClient("http://localhost:8081")
	if client.httpClient.Timeout != 10*time.Second {
		t.Errorf("Expected timeout 10s, got %v", client.httpClient.Timeout)
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

func TestCreatePlanBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	plan := PlanRequest{Title: "Тест", Date: "2024-12-25"}

	err := client.CreatePlan(plan)
	if err == nil {
		t.Error("Expected error for bad request")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Error should contain status code 400, got: %v", err)
	}
}

func TestCreatePlanNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	plan := PlanRequest{Title: "Тест", Date: "2024-12-25"}

	err := client.CreatePlan(plan)
	if err == nil {
		t.Error("Expected error for not found")
	}
}

func TestCreatePlanRequestFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		var plan PlanRequest
		json.NewDecoder(r.Body).Decode(&plan)

		if plan.Title != "Тест" {
			t.Errorf("Expected title 'Тест', got '%s'", plan.Title)
		}
		if plan.Description != "Описание" {
			t.Errorf("Expected description 'Описание', got '%s'", plan.Description)
		}
		if plan.Date != "2024-12-25" {
			t.Errorf("Expected date '2024-12-25', got '%s'", plan.Date)
		}
		if plan.Time != "14:00" {
			t.Errorf("Expected time '14:00', got '%s'", plan.Time)
		}
		if plan.IsAllDay != false {
			t.Errorf("Expected is_all_day false, got true")
		}

		w.WriteHeader(http.StatusCreated)
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

func TestCreatePlanAllDay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var plan PlanRequest
		json.NewDecoder(r.Body).Decode(&plan)

		if plan.IsAllDay != true {
			t.Errorf("Expected is_all_day true, got false")
		}
		if plan.Time != "" {
			t.Errorf("Expected empty time, got '%s'", plan.Time)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	plan := PlanRequest{
		Title:       "Дедлайн",
		Description: "Сдать отчёт",
		Date:        "2024-12-26",
		Time:        "",
		IsAllDay:    true,
	}

	err := client.CreatePlan(plan)
	if err != nil {
		t.Errorf("CreatePlan returned error: %v", err)
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

func TestGetPlansByDateEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetPlansByDate("2099-01-01")

	if err != nil {
		t.Errorf("GetPlansByDate returned error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 plans, got %d", len(result))
	}
}

func TestGetPlansByDateInvalidDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetPlansByDate("invalid-date")

	if err != nil {
		t.Errorf("GetPlansByDate returned error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 plans for invalid date, got %d", len(result))
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

func TestGetPlansByDateRangeEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetPlansByDateRange("2099-01-01", "2099-01-31")

	if err != nil {
		t.Errorf("GetPlansByDateRange returned error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 plans, got %d", len(result))
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

func TestGetPlansByDateRangeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPlansByDateRange("2024-12-23", "2024-12-29")

	if err == nil {
		t.Error("Expected error from GetPlansByDateRange")
	}
}

func TestGetPlansByDateRangeBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing parameters"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPlansByDateRange("2024-12-23", "2024-12-29")

	if err == nil {
		t.Error("Expected error for bad request")
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

func TestGetJSONEmptyArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var result []Plan
	err := client.getJSON("/plans", &result)

	if err != nil {
		t.Errorf("Unexpected error for empty array: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d", len(result))
	}
}

func TestGetJSONServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var result []Plan
	err := client.getJSON("/plans", &result)

	if err == nil {
		t.Error("Expected error for server error")
	}
}

func TestGetJSONNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var result []Plan
	err := client.getJSON("/plans", &result)

	if err == nil {
		t.Error("Expected error for not found")
	}
}

func TestClientConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1")
	err := client.CreatePlan(PlanRequest{Title: "Test", Date: "2024-12-25"})

	if err == nil {
		t.Error("Expected connection error")
	}
}

func TestClientGetPlansConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1")
	_, err := client.GetPlansByDate("2024-12-25")

	if err == nil {
		t.Error("Expected connection error")
	}
}

func TestClientGetPlansRangeConnectionError(t *testing.T) {
	client := NewClient("http://localhost:1")
	_, err := client.GetPlansByDateRange("2024-12-23", "2024-12-29")

	if err == nil {
		t.Error("Expected connection error")
	}
}

func TestCreatePlanTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	client.httpClient.Timeout = 50 * time.Millisecond

	err := client.CreatePlan(PlanRequest{Title: "Test", Date: "2024-12-25"})
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestGetPlansByDateTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	client.httpClient.Timeout = 50 * time.Millisecond

	_, err := client.GetPlansByDate("2024-12-25")
	if err == nil {
		t.Error("Expected timeout error")
	}
}

func TestCreatePlanCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/plans" {
			t.Errorf("Expected path /plans, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.CreatePlan(PlanRequest{Title: "Test", Date: "2024-12-25"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetPlansByDateCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/plans" {
			t.Errorf("Expected path /plans, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPlansByDate("2024-12-25")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetPlansByDateRangeCorrectPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/plans/range" {
			t.Errorf("Expected path /plans/range, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPlansByDateRange("2024-12-23", "2024-12-29")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetPlansByDateQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		if date != "2024-12-25" {
			t.Errorf("Expected date=2024-12-25, got %s", date)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPlansByDate("2024-12-25")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetPlansByDateRangeQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		if start != "2024-12-01" {
			t.Errorf("Expected start=2024-12-01, got %s", start)
		}
		if end != "2024-12-31" {
			t.Errorf("Expected end=2024-12-31, got %s", end)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetPlansByDateRange("2024-12-01", "2024-12-31")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetPlansByDateReturnsFullPlanData(t *testing.T) {
	plans := []Plan{
		{
			ID:          1,
			Title:       "Тест",
			Description: "Описание",
			Date:        "2024-12-25",
			Time:        "14:00",
			IsAllDay:    false,
			CreatedAt:   "2024-12-24T10:00:00Z",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.GetPlansByDate("2024-12-25")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 plan, got %d", len(result))
	}

	plan := result[0]
	if plan.ID != 1 {
		t.Errorf("Expected ID 1, got %d", plan.ID)
	}
	if plan.Title != "Тест" {
		t.Errorf("Expected title 'Тест', got '%s'", plan.Title)
	}
	if plan.Description != "Описание" {
		t.Errorf("Expected description 'Описание', got '%s'", plan.Description)
	}
	if plan.Date != "2024-12-25" {
		t.Errorf("Expected date '2024-12-25', got '%s'", plan.Date)
	}
	if plan.Time != "14:00" {
		t.Errorf("Expected time '14:00', got '%s'", plan.Time)
	}
	if plan.IsAllDay != false {
		t.Errorf("Expected is_all_day false, got true")
	}
	if plan.CreatedAt != "2024-12-24T10:00:00Z" {
		t.Errorf("Expected created_at '2024-12-24T10:00:00Z', got '%s'", plan.CreatedAt)
	}
}
