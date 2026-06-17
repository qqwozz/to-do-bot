package bot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClientInitialized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	if client.baseURL != server.URL {
		t.Errorf("baseURL = %s, want %s", client.baseURL, server.URL)
	}
}

func TestPlanSerialization(t *testing.T) {
	plan := Plan{
		ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25",
		Time: "14:00", IsAllDay: false, CreatedAt: "2024-12-24T10:00:00Z",
	}
	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded Plan
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.ID != plan.ID || decoded.Title != plan.Title || decoded.CreatedAt != plan.CreatedAt {
		t.Errorf("mismatch: %+v != %+v", decoded, plan)
	}
}

func TestPlanRequestJSONTags(t *testing.T) {
	plan := PlanRequest{Title: "Заголовок", Description: "Описание", Date: "2025-01-15", Time: "09:30", IsAllDay: true}
	data, _ := json.Marshal(plan)
	s := string(data)
	for _, field := range []string{`"title"`, `"description"`, `"date"`, `"time"`, `"is_all_day"`} {
		if !strings.Contains(s, field) {
			t.Errorf("JSON missing field %s", field)
		}
	}
}
