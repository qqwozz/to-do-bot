package bot

import (
	"testing"
)

func TestMainMenuKeyboardLayout(t *testing.T) {
	kb := mainMenuKeyboard()
	if len(kb.InlineKeyboard) != 1 {
		t.Fatalf("expected 1 row, got %d", len(kb.InlineKeyboard))
	}
	if len(kb.InlineKeyboard[0]) != 2 {
		t.Fatalf("expected 2 buttons in row, got %d", len(kb.InlineKeyboard[0]))
	}
}

func TestMainMenuKeyboardTexts(t *testing.T) {
	kb := mainMenuKeyboard()
	expected := []struct {
		text string
		data string
	}{
		{"➕ Создать план", "plan:create"},
		{"📋 Мои планы", "plans:show"},
	}
	for i, exp := range expected {
		btn := kb.InlineKeyboard[0][i]
		if btn.Text != exp.text {
			t.Errorf("button %d: text = %q, want %q", i, btn.Text, exp.text)
		}
		if *btn.CallbackData != exp.data {
			t.Errorf("button %d: data = %q, want %q", i, *btn.CallbackData, exp.data)
		}
	}
}

func TestPlansPeriodKeyboardLayout(t *testing.T) {
	kb := plansPeriodKeyboard()
	if len(kb.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(kb.InlineKeyboard))
	}
}

func TestPlansPeriodKeyboardCallbacks(t *testing.T) {
	kb := plansPeriodKeyboard()
	expected := map[string]string{
		"plans:today":     "📅 Сегодня",
		"plans:tomorrow":  "📅 Завтра",
		"plans:week":      "📆 Эта неделя",
		"plans:nextweek":  "📆 След. неделя",
		"menu:main":       "◀️ Назад",
	}

	found := make(map[string]bool)
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData == nil {
				continue
			}
			data := *btn.CallbackData
			if text, ok := expected[data]; ok {
				if btn.Text != text {
					t.Errorf("callback %s: text = %q, want %q", data, btn.Text, text)
				}
				found[data] = true
			}
		}
	}
	for data := range expected {
		if !found[data] {
			t.Errorf("missing callback: %s", data)
		}
	}
}

func TestPlansViewKeyboardLayout(t *testing.T) {
	kb := plansViewKeyboard()
	if len(kb.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(kb.InlineKeyboard))
	}
	total := 0
	for _, row := range kb.InlineKeyboard {
		total += len(row)
	}
	if total != 6 {
		t.Fatalf("expected 6 total buttons, got %d", total)
	}
}

func TestPlansViewKeyboardHasBackButton(t *testing.T) {
	kb := plansViewKeyboard()
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil && *btn.CallbackData == "menu:main" {
				if btn.Text != "◀️ Назад" {
					t.Errorf("back button text = %q, want '◀️ Назад'", btn.Text)
				}
				return
			}
		}
	}
	t.Error("missing back button")
}

func TestCreateDateKeyboardLayout(t *testing.T) {
	kb := createDateKeyboard()
	if len(kb.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(kb.InlineKeyboard))
	}
}

func TestCreateDateKeyboardCallbacks(t *testing.T) {
	kb := createDateKeyboard()
	expected := map[string]bool{
		"plan:date:today":    true,
		"plan:date:tomorrow": true,
		"plan:date:week":     true,
		"plan:date:custom":   true,
		"plan:cancel":        true,
	}

	found := make(map[string]bool)
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil {
				found[*btn.CallbackData] = true
			}
		}
	}
	for data := range expected {
		if !found[data] {
			t.Errorf("missing callback: %s", data)
		}
	}
}

func TestCreateTimeKeyboardLayout(t *testing.T) {
	kb := createTimeKeyboard()
	if len(kb.InlineKeyboard) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(kb.InlineKeyboard))
	}
}

func TestCreateTimeKeyboardCallbacks(t *testing.T) {
	kb := createTimeKeyboard()
	expected := map[string]bool{
		"plan:time:morning":   true,
		"plan:time:afternoon": true,
		"plan:time:evening":   true,
		"plan:time:allday":    true,
		"plan:cancel":         true,
	}

	found := make(map[string]bool)
	for _, row := range kb.InlineKeyboard {
		for _, btn := range row {
			if btn.CallbackData != nil {
				found[*btn.CallbackData] = true
			}
		}
	}
	for data := range expected {
		if !found[data] {
			t.Errorf("missing callback: %s", data)
		}
	}
}

func TestCancelKeyboardLayout(t *testing.T) {
	kb := cancelKeyboard()
	if len(kb.InlineKeyboard) != 1 {
		t.Fatalf("expected 1 row, got %d", len(kb.InlineKeyboard))
	}
	if len(kb.InlineKeyboard[0]) != 1 {
		t.Fatalf("expected 1 button, got %d", len(kb.InlineKeyboard[0]))
	}
	btn := kb.InlineKeyboard[0][0]
	if *btn.CallbackData != "plan:cancel" {
		t.Errorf("callback = %q, want 'plan:cancel'", *btn.CallbackData)
	}
	if btn.Text != "❌ Отмена" {
		t.Errorf("text = %q, want '❌ Отмена'", btn.Text)
	}
}
