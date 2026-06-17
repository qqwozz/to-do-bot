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
	c := NewClient("http://localhost:8081")
	if c.baseURL != "http://localhost:8081" {
		t.Errorf("baseURL = %s", c.baseURL)
	}
	if c.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

func TestNewClientTimeout(t *testing.T) {
	c := NewClient("http://localhost:8081")
	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("timeout = %v", c.httpClient.Timeout)
	}
}

func TestCreatePlanSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/plans" {
			t.Errorf("unexpected: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.CreatePlan(PlanRequest{Title: "Тест", Date: "2024-12-25"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreatePlanError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.CreatePlan(PlanRequest{Title: "Тест", Date: "2024-12-25"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestCreatePlanBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad"))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.CreatePlan(PlanRequest{Title: "Тест", Date: "2024-12-25"})
	if err == nil || !strings.Contains(err.Error(), "400") {
		t.Errorf("expected 400 error, got: %v", err)
	}
}

func TestCreatePlanNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.CreatePlan(PlanRequest{Title: "Тест", Date: "2024-12-25"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestCreatePlanContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %s", ct)
		}
		var plan PlanRequest
		json.NewDecoder(r.Body).Decode(&plan)
		if plan.Title != "Тест" || plan.Date != "2024-12-25" || plan.Time != "14:00" {
			t.Errorf("plan fields mismatch: %+v", plan)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.CreatePlan(PlanRequest{Title: "Тест", Date: "2024-12-25", Time: "14:00"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCreatePlanAllDay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var plan PlanRequest
		json.NewDecoder(r.Body).Decode(&plan)
		if !plan.IsAllDay {
			t.Error("expected is_all_day=true")
		}
		if plan.Time != "" {
			t.Errorf("expected empty time, got %s", plan.Time)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.CreatePlan(PlanRequest{Title: "Дедлайн", Date: "2024-12-26", IsAllDay: true})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetPlansByDate(t *testing.T) {
	plans := []Plan{
		{ID: 1, Title: "План 1", Date: "2024-12-25"},
		{ID: 2, Title: "План 2", Date: "2024-12-25"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("date") != "2024-12-25" {
			t.Error("wrong date param")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	result, err := c.GetPlansByDate("2024-12-25")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 plans, got %d", len(result))
	}
}

func TestGetPlansByDateEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	result, err := c.GetPlansByDate("2099-01-01")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 plans, got %d", len(result))
	}
}

func TestGetPlansByDateError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.GetPlansByDate("2024-12-25")
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetPlansByDateRange(t *testing.T) {
	plans := []Plan{
		{ID: 1, Title: "План 1", Date: "2024-12-23"},
		{ID: 2, Title: "План 2", Date: "2024-12-25"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start, end := r.URL.Query().Get("start"), r.URL.Query().Get("end")
		if start != "2024-12-23" || end != "2024-12-29" {
			t.Errorf("wrong params: start=%s end=%s", start, end)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	result, err := c.GetPlansByDateRange("2024-12-23", "2024-12-29")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 plans, got %d", len(result))
	}
}

func TestGetPlansByDateRangeEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Plan{})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	result, err := c.GetPlansByDateRange("2099-01-01", "2099-01-31")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 plans, got %d", len(result))
	}
}

func TestGetPlansByDateRangeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	_, err := c.GetPlansByDateRange("2024-12-23", "2024-12-29")
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetJSONInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid"))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	var result []Plan
	err := c.getJSON("/plans", &result)
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetJSONEmptyArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	c := NewClient(server.URL)
	var result []Plan
	err := c.getJSON("/plans", &result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}

func TestClientConnectionError(t *testing.T) {
	c := NewClient("http://localhost:1")
	err := c.CreatePlan(PlanRequest{Title: "Test", Date: "2024-12-25"})
	if err == nil {
		t.Error("expected connection error")
	}
}

func TestClientGetPlansConnectionError(t *testing.T) {
	c := NewClient("http://localhost:1")
	_, err := c.GetPlansByDate("2024-12-25")
	if err == nil {
		t.Error("expected connection error")
	}
}

func TestClientGetPlansRangeConnectionError(t *testing.T) {
	c := NewClient("http://localhost:1")
	_, err := c.GetPlansByDateRange("2024-12-23", "2024-12-29")
	if err == nil {
		t.Error("expected connection error")
	}
}

func TestGetPlansByDateReturnsFullData(t *testing.T) {
	plans := []Plan{
		{ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25", Time: "14:00", IsAllDay: false, CreatedAt: "2024-12-24T10:00:00Z"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plans)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	result, err := c.GetPlansByDate("2024-12-25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(result))
	}

	p := result[0]
	if p.ID != 1 || p.Title != "Тест" || p.Description != "Описание" || p.Date != "2024-12-25" || p.Time != "14:00" || p.IsAllDay || p.CreatedAt != "2024-12-24T10:00:00Z" {
		t.Errorf("plan fields mismatch: %+v", p)
	}
}
