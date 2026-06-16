package bot

import (
	"testing"
)

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2024-12-25", "25.12.2024"},
		{"2024-01-01", "01.01.2024"},
		{"invalid-date", "invalid-date"},
		{"", ""},
	}

	for _, tt := range tests {
		result := formatDate(tt.input)
		if result != tt.expected {
			t.Errorf("formatDate(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatPlansListEmpty(t *testing.T) {
	b := &Bot{}
	result := b.formatPlansList([]Plan{}, "Сегодня", "", "")

	if len(result) == 0 {
		t.Error("formatPlansList returned empty string")
	}
	if !contains(result, "Нет планов") {
		t.Error("Expected 'Нет планов' for empty list")
	}
}

func TestFormatPlansListWithPlans(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{
			ID:          1,
			Title:       "Встреча",
			Description: "Обсудить проект",
			Date:        "2024-12-25",
			Time:        "14:00",
			IsAllDay:    false,
		},
		{
			ID:          2,
			Title:       "Дедлайн",
			Description: "Сдать отчёт",
			Date:        "2024-12-26",
			IsAllDay:    true,
		},
	}

	result := b.formatPlansList(plans, "Завтра", "", "")

	if !contains(result, "Встреча") {
		t.Error("Expected plan title in output")
	}
	if !contains(result, "14:00") {
		t.Error("Expected time in output")
	}
	if !contains(result, "Весь день") {
		t.Error("Expected 'Весь день' for all-day plan")
	}
}

func TestFormatPlansListWithSubtitle(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25"},
	}

	result := b.formatPlansList(plans, "Неделя", "23.12 - 29.12", "")

	if !contains(result, "23.12 - 29.12") {
		t.Error("Expected subtitle in output")
	}
}

func TestCreateNavigationKeyboard(t *testing.T) {
	b := &Bot{}
	keyboard := b.createNavigationKeyboard()

	if len(keyboard.InlineKeyboard) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(keyboard.InlineKeyboard))
	}
	if len(keyboard.InlineKeyboard[0]) != 2 {
		t.Errorf("Expected 2 buttons in first row, got %d", len(keyboard.InlineKeyboard[0]))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
