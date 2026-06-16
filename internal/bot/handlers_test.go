package bot

import (
	"fmt"
	"strings"
	"testing"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func containsSubstr(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2024-12-25", "25.12.2024"},
		{"2024-01-01", "01.01.2024"},
		{"2025-06-15", "15.06.2025"},
		{"2000-02-29", "29.02.2000"},
		{"invalid-date", "invalid-date"},
		{"", ""},
		{"2024/12/25", "2024/12/25"},
		{"25-12-2024", "25-12-2024"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatDate(tt.input)
			if result != tt.expected {
				t.Errorf("formatDate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
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
	if !contains(result, "Сегодня") {
		t.Error("Expected title 'Сегодня' in output")
	}
}

func TestFormatPlansListNilPlans(t *testing.T) {
	b := &Bot{}
	result := b.formatPlansList(nil, "Завтра", "", "")

	if !contains(result, "Нет планов") {
		t.Error("Expected 'Нет планов' for nil plans")
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
	if !contains(result, "Дедлайн") {
		t.Error("Expected second plan title")
	}
	if !contains(result, "Обсудить проект") {
		t.Error("Expected first plan description")
	}
	if !contains(result, "Сдать отчёт") {
		t.Error("Expected second plan description")
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

func TestFormatPlansListWithoutSubtitle(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25"},
	}

	result := b.formatPlansList(plans, "Сегодня", "", "")

	if contains(result, "23.12 - 29.12") {
		t.Error("Should not contain subtitle when empty")
	}
}

func TestFormatPlansListMultiplePlansNumbering(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "План 1", Description: "Описание 1", Date: "2024-12-25", Time: "09:00", IsAllDay: false},
		{ID: 2, Title: "План 2", Description: "Описание 2", Date: "2024-12-25", Time: "10:00", IsAllDay: false},
		{ID: 3, Title: "План 3", Description: "Описание 3", Date: "2024-12-25", Time: "11:00", IsAllDay: false},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	if !contains(result, "*1.*") {
		t.Error("Expected numbering '*1.*'")
	}
	if !contains(result, "*2.*") {
		t.Error("Expected numbering '*2.*'")
	}
	if !contains(result, "*3.*") {
		t.Error("Expected numbering '*3.*'")
	}
}

func TestFormatPlansListTimeWithEmptyTime(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25", Time: "", IsAllDay: false},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	if !contains(result, "Весь день") {
		t.Error("Empty time with IsAllDay=false should show 'Весь день'")
	}
}

func TestFormatPlansListTimeWithNonEmptyTime(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "Тест", Description: "Описание", Date: "2024-12-25", Time: "15:30", IsAllDay: false},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	if !contains(result, "15:30") {
		t.Error("Should contain the time")
	}
	if contains(result, "Весь день") {
		t.Error("Should not show 'Весь день' when time is set")
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
	if len(keyboard.InlineKeyboard[1]) != 2 {
		t.Errorf("Expected 2 buttons in second row, got %d", len(keyboard.InlineKeyboard[1]))
	}
}

func TestCreateNavigationKeyboardButtons(t *testing.T) {
	b := &Bot{}
	keyboard := b.createNavigationKeyboard()

	expectedButtons := map[string]string{
		"today":    "📅 Сегодня",
		"tomorrow": "📅 Завтра",
		"week":     "📆 Эта неделя",
		"nextweek": "📆 След. неделя",
	}

	found := make(map[string]bool)
	for _, row := range keyboard.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData == nil {
				continue
			}
			if expected, ok := expectedButtons[*btn.CallbackData]; ok {
				if btn.Text != expected {
					t.Errorf("Button %s: expected text %q, got %q", *btn.CallbackData, expected, btn.Text)
				}
				found[*btn.CallbackData] = true
			}
		}
	}

	for key, expected := range expectedButtons {
		if !found[key] {
			t.Errorf("Missing button with callback data %q and text %q", key, expected)
		}
	}
}

func TestParseCreatePlanInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		title       string
		description string
		date        string
		time        string
		isAllDay    bool
		valid       bool
	}{
		{
			name:        "full plan with time",
			input:       "CREATE:Встреча|Обсудить проект|2024-12-25|14:00",
			title:       "Встреча",
			description: "Обсудить проект",
			date:        "2024-12-25",
			time:        "14:00",
			isAllDay:    false,
			valid:       true,
		},
		{
			name:        "plan without time (all day)",
			input:       "CREATE:Дедлайн|Сдать отчёт|2024-12-26|",
			title:       "Дедлайн",
			description: "Сдать отчёт",
			date:        "2024-12-26",
			time:        "",
			isAllDay:    true,
			valid:       true,
		},
		{
			name:        "plan without time field",
			input:       "CREATE:Тест|Описание|2024-12-25",
			title:       "Тест",
			description: "Описание",
			date:        "2024-12-25",
			time:        "",
			isAllDay:    true,
			valid:       true,
		},
		{
			name:  "too few parts",
			input: "CREATE:Тест|Описание",
			valid: false,
		},
		{
			name:  "only one part",
			input: "CREATE:Тест",
			valid: false,
		},
		{
			name:  "empty create",
			input: "CREATE:",
			valid: false,
		},
		{
			name:  "no create prefix",
			input: "Hello",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(strings.TrimPrefix(tt.input, "CREATE:"), "|")

			if tt.valid {
				if len(parts) < 3 {
					t.Errorf("Expected valid input, got %d parts", len(parts))
					return
				}

				title := strings.TrimSpace(parts[0])
				description := strings.TrimSpace(parts[1])
				date := strings.TrimSpace(parts[2])

				if title != tt.title {
					t.Errorf("Title: got %q, want %q", title, tt.title)
				}
				if description != tt.description {
					t.Errorf("Description: got %q, want %q", description, tt.description)
				}
				if date != tt.date {
					t.Errorf("Date: got %q, want %q", date, tt.date)
				}

				time := ""
				isAllDay := true
				if len(parts) > 3 && strings.TrimSpace(parts[3]) != "" {
					time = strings.TrimSpace(parts[3])
					isAllDay = false
				}

				if time != tt.time {
					t.Errorf("Time: got %q, want %q", time, tt.time)
				}
				if isAllDay != tt.isAllDay {
					t.Errorf("IsAllDay: got %v, want %v", isAllDay, tt.isAllDay)
				}
			} else {
				if len(parts) >= 3 {
					t.Error("Expected invalid input, but got enough parts")
				}
			}
		})
	}
}

func TestFormatPlansListTitle(t *testing.T) {
	b := &Bot{}

	result := b.formatPlansList([]Plan{}, "Мой заголовок", "", "")

	if !contains(result, "Мой заголовок") {
		t.Error("Expected custom title in output")
	}
}

func TestFormatPlansListFooterIgnored(t *testing.T) {
	b := &Bot{}

	result := b.formatPlansList([]Plan{}, "Тест", "", "footer text")

	if contains(result, "footer text") {
		t.Error("Footer should not appear in output")
	}
}

func TestFormatPlansListPlanOrdering(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 2, Title: "Второй", Description: "Desc2", Date: "2024-12-25", Time: "10:00", IsAllDay: false},
		{ID: 1, Title: "Первый", Description: "Desc1", Date: "2024-12-25", Time: "09:00", IsAllDay: false},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	idx1 := strings.Index(result, "Второй")
	idx2 := strings.Index(result, "Первый")

	if idx1 >= idx2 {
		t.Error("Plans should be in the order provided")
	}
}

func TestContainsHelper(t *testing.T) {
	tests := []struct {
		s, substr string
		expected  bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", true},
		{"hello world", "xyz", false},
		{"", "", true},
		{"hello", "", true},
		{"", "hello", false},
		{"abc", "abcd", false},
	}

	for _, tt := range tests {
		result := contains(tt.s, tt.substr)
		if result != tt.expected {
			t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
		}
	}
}

func TestContainsSubstrHelper(t *testing.T) {
	tests := []struct {
		s, substr string
		expected  bool
	}{
		{"hello world", "world", true},
		{"hello world", "xyz", false},
		{"a", "a", true},
		{"ab", "abc", false},
	}

	for _, tt := range tests {
		result := containsSubstr(tt.s, tt.substr)
		if result != tt.expected {
			t.Errorf("containsSubstr(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
		}
	}
}

func TestFormatPlansListDescription(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "Тест", Description: "Очень длинное описание дела", Date: "2024-12-25"},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	if !contains(result, "Очень длинное описание дела") {
		t.Error("Expected full description in output")
	}
}

func TestFormatPlansListMultipleDescriptions(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "План 1", Description: "Описание 1", Date: "2024-12-25"},
		{ID: 2, Title: "План 2", Description: "Описание 2", Date: "2024-12-25"},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	if !contains(result, "Описание 1") {
		t.Error("Expected first description")
	}
	if !contains(result, "Описание 2") {
		t.Error("Expected second description")
	}
}

func TestFormatPlansListEmptyDescription(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{ID: 1, Title: "Тест", Description: "", Date: "2024-12-25"},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	if !contains(result, "Тест") {
		t.Error("Expected title in output")
	}
}

func TestFormatPlansListWithAllFields(t *testing.T) {
	b := &Bot{}
	plans := []Plan{
		{
			ID:          42,
			Title:       "Важная встреча",
			Description: "Обсудить квартальные результаты",
			Date:        "2024-12-25",
			Time:        "14:00",
			IsAllDay:    false,
			CreatedAt:   "2024-12-20T10:00:00Z",
		},
	}

	result := b.formatPlansList(plans, "Тест", "", "")

	if !contains(result, "Важная встреча") {
		t.Error("Expected title")
	}
	if !contains(result, "Обсудить квартальные результаты") {
		t.Error("Expected description")
	}
	if !contains(result, "14:00") {
		t.Error("Expected time")
	}
	if contains(result, "42") {
		t.Error("Should not contain ID")
	}
	if contains(result, "2024-12-20") {
		t.Error("Should not contain CreatedAt")
	}
}

func TestCreateNavigationKeyboardRowCount(t *testing.T) {
	b := &Bot{}
	keyboard := b.createNavigationKeyboard()

	if len(keyboard.InlineKeyboard) != 2 {
		t.Errorf("Expected exactly 2 rows, got %d", len(keyboard.InlineKeyboard))
	}
}

func TestCreateNavigationKeyboardTotalButtons(t *testing.T) {
	b := &Bot{}
	keyboard := b.createNavigationKeyboard()

	total := 0
	for _, row := range keyboard.InlineKeyboard {
		total += len(row)
	}

	if total != 4 {
		t.Errorf("Expected 4 total buttons, got %d", total)
	}
}

func TestCreateNavigationKeyboardCallbackData(t *testing.T) {
	b := &Bot{}
	keyboard := b.createNavigationKeyboard()

	expectedCallbacks := []string{"today", "tomorrow", "week", "nextweek"}
	foundCallbacks := make(map[string]bool)

	for _, row := range keyboard.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil {
				foundCallbacks[*btn.CallbackData] = true
			}
		}
	}

	for _, cb := range expectedCallbacks {
		if !foundCallbacks[cb] {
			t.Errorf("Missing callback data: %s", cb)
		}
	}
}

func TestFormatDateConsistency(t *testing.T) {
	for month := 1; month <= 12; month++ {
		date := fmt.Sprintf("2024-%02d-01", month)
		result := formatDate(date)
		if result == date {
			t.Errorf("formatDate(%q) returned original date, expected formatted", date)
		}
	}
}

func TestFormatDateEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"leap year", "2024-02-29", "29.02.2024"},
		{"year 2000", "2000-01-01", "01.01.2000"},
		{"year 1999", "1999-12-31", "31.12.1999"},
		{"single digit month", "2024-1-15", "2024-1-15"},
		{"wrong separator", "2024.12.25", "2024.12.25"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDate(tt.input)
			if result != tt.expected {
				t.Errorf("formatDate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
