package bot

import (
	"testing"
)

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"2024-12-25", "25.12.2024"},
		{"2024-01-01", "01.01.2024"},
		{"2025-06-15", "15.06.2025"},
		{"2000-02-29", "29.02.2000"},
		{"invalid", "invalid"},
		{"", ""},
		{"2024/12/25", "2024/12/25"},
		{"25-12-2024", "25-12-2024"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := formatDate(tt.input); got != tt.expected {
				t.Errorf("formatDate(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatPlansListEmpty(t *testing.T) {
	result := formatPlansList([]Plan{}, "Сегодня", "")
	if result == "" {
		t.Error("empty result")
	}
	if !containsStr(result, "Нет планов") {
		t.Error("missing 'Нет планов'")
	}
}

func TestFormatPlansListNil(t *testing.T) {
	result := formatPlansList(nil, "Завтра", "")
	if !containsStr(result, "Нет планов") {
		t.Error("nil plans should show 'Нет планов'")
	}
}

func TestFormatPlansListWithPlans(t *testing.T) {
	plans := []Plan{
		{ID: 1, Title: "Встреча", Description: "Обсудить проект", Date: "2024-12-25", Time: "14:00", IsAllDay: false},
		{ID: 2, Title: "Дедлайн", Description: "Сдать отчёт", Date: "2024-12-26", IsAllDay: true},
	}
	result := formatPlansList(plans, "Завтра", "")

	for _, s := range []string{"Встреча", "14:00", "Дедлайн", "Весь день", "Обсудить проект", "Сдать отчёт"} {
		if !containsStr(result, s) {
			t.Errorf("missing %q", s)
		}
	}
}

func TestFormatPlansListSubtitle(t *testing.T) {
	plans := []Plan{{ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25"}}
	result := formatPlansList(plans, "Неделя", "23.12 — 29.12")
	if !containsStr(result, "23.12 — 29.12") {
		t.Error("missing subtitle")
	}
}

func TestFormatPlansListNoSubtitle(t *testing.T) {
	plans := []Plan{{ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25"}}
	result := formatPlansList(plans, "Сегодня", "")
	if containsStr(result, "23.12") {
		t.Error("should not contain subtitle when empty")
	}
}

func TestFormatPlansListNumbering(t *testing.T) {
	plans := []Plan{
		{ID: 1, Title: "A", Description: "a", Date: "2024-12-25", Time: "09:00", IsAllDay: false},
		{ID: 2, Title: "B", Description: "b", Date: "2024-12-25", Time: "10:00", IsAllDay: false},
		{ID: 3, Title: "C", Description: "c", Date: "2024-12-25", Time: "11:00", IsAllDay: false},
	}
	result := formatPlansList(plans, "Тест", "")
	for _, n := range []string{"*1.*", "*2.*", "*3.*"} {
		if !containsStr(result, n) {
			t.Errorf("missing numbering %s", n)
		}
	}
}

func TestFormatPlansListTimeVsAllDay(t *testing.T) {
	timed := []Plan{{ID: 1, Title: "T", Description: "D", Date: "2024-12-25", Time: "15:30", IsAllDay: false}}
	allday := []Plan{{ID: 1, Title: "T", Description: "D", Date: "2024-12-25", IsAllDay: true}}

	r1 := formatPlansList(timed, "T", "")
	if !containsStr(r1, "15:30") {
		t.Error("missing time")
	}
	if containsStr(r1, "Весь день") {
		t.Error("should not show 'Весь день' when time is set")
	}

	r2 := formatPlansList(allday, "T", "")
	if !containsStr(r2, "Весь день") {
		t.Error("missing 'Весь день'")
	}
}

func TestFormatPlansListTitle(t *testing.T) {
	result := formatPlansList([]Plan{}, "Мой заголовок", "")
	if !containsStr(result, "Мой заголовок") {
		t.Error("missing custom title")
	}
}

func TestFormatDateConsistency(t *testing.T) {
	for month := 1; month <= 12; month++ {
		date := "2024-" + pad(month) + "-01"
		if got := formatDate(date); got == date {
			t.Errorf("formatDate(%q) not formatted", date)
		}
	}
}

func pad(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return "1" + string(rune('0'+n-10))
}

func containsStr(s, sub string) bool {
	return len(s) > 0 && len(sub) > 0 && (s == sub || len(s) > len(sub) && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
