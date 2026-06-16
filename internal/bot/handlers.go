package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	switch {
	case text == "/start":
		b.sendMainMenu(chatID)
	case text == "/today":
		b.sendPlansForDate(chatID, time.Now().Format("2006-01-02"), "Сегодня")
	case text == "/tomorrow":
		tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		b.sendPlansForDate(chatID, tomorrow, "Завтра")
	case text == "/week":
		b.sendPlansForWeek(chatID, 0, "Эта неделя")
	case text == "/nextweek":
		b.sendPlansForWeek(chatID, 7, "Следующая неделя")
	case strings.HasPrefix(text, "CREATE:"):
		b.handleCreatePlan(chatID, text)
	default:
		b.SendMessage(chatID, "Используйте /start для начала работы.")
	}
}

func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID

	switch query.Data {
	case "create_plan":
		b.sendCreatePlanPrompt(chatID)
	case "view_plans":
		b.sendMainMenu(chatID)
	case "today":
		b.sendPlansForDate(chatID, time.Now().Format("2006-01-02"), "Сегодня")
	case "tomorrow":
		tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		b.sendPlansForDate(chatID, tomorrow, "Завтра")
	case "week":
		b.sendPlansForWeek(chatID, 0, "Эта неделя")
	case "nextweek":
		b.sendPlansForWeek(chatID, 7, "Следующая неделя")
	}

	b.AnswerCallback(query.ID)
}

func (b *Bot) sendMainMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать план", "create_plan"),
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои планы", "view_plans"),
		),
	)

	text := "👋 Добро пожаловать в To-Do бот!\n\n" +
		"Я помогу вам управлять списком задач.\n\n" +
		"Выберите действие:"

	b.SendMessageWithKeyboard(chatID, text, keyboard)
}

func (b *Bot) sendCreatePlanPrompt(chatID int64) {
	text := "📝 *Создание нового плана*\n\n" +
		"Отправьте план в формате:\n" +
		"`Название|Описание|Дата|Время`\n\n" +
		"*Примеры:*\n" +
		"• `Встреча|Обсудить проект|2024-12-25|14:00`\n" +
		"• `Дедлайн|Сдать отчёт|2024-12-26|` (весь день)\n\n" +
		"📅 Формат даты: ГГГГ-ММ-ДД\n" +
		"🕐 Формат времени: ЧЧ:ММ (или оставьте пустым)"

	b.SendMessage(chatID, text)
}

func (b *Bot) handleCreatePlan(chatID int64, text string) {
	parts := strings.Split(strings.TrimPrefix(text, "CREATE:"), "|")

	if len(parts) < 3 {
		b.SendMessage(chatID, "❌ Неверный формат. Используйте: Название|Описание|Дата")
		return
	}

	plan := PlanRequest{
		Title:       strings.TrimSpace(parts[0]),
		Description: strings.TrimSpace(parts[1]),
		Date:        strings.TrimSpace(parts[2]),
	}

	if len(parts) > 3 && strings.TrimSpace(parts[3]) != "" {
		plan.Time = strings.TrimSpace(parts[3])
		plan.IsAllDay = false
	} else {
		plan.IsAllDay = true
	}

	if err := b.client.CreatePlan(plan); err != nil {
		b.SendMessage(chatID, fmt.Sprintf("❌ Ошибка: %v", err))
		return
	}

	b.SendMessage(chatID, "✅ План успешно создан!")
}

func (b *Bot) sendPlansForDate(chatID int64, date string, label string) {
	plans, err := b.client.GetPlansByDate(date)
	if err != nil {
		b.SendMessage(chatID, "❌ Ошибка получения планов")
		return
	}

	text := b.formatPlansList(plans, label, "", "")
	keyboard := b.createNavigationKeyboard()
	b.SendMessageWithKeyboard(chatID, text, keyboard)
}

func (b *Bot) sendPlansForWeek(chatID int64, offsetDays int, label string) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+offsetDays)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	startDate := startOfWeek.Format("2006-01-02")
	endDate := endOfWeek.Format("2006-01-02")

	plans, err := b.client.GetPlansByDateRange(startDate, endDate)
	if err != nil {
		b.SendMessage(chatID, "❌ Ошибка получения планов")
		return
	}

	subtitle := fmt.Sprintf("%s - %s", formatDate(startDate), formatDate(endDate))
	text := b.formatPlansList(plans, label, subtitle, "")
	keyboard := b.createNavigationKeyboard()
	b.SendMessageWithKeyboard(chatID, text, keyboard)
}

func (b *Bot) formatPlansList(plans []Plan, title string, subtitle string, footer string) string {
	if len(plans) == 0 {
		return fmt.Sprintf("📋 *%s*\n\nНет планов.", title)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 *%s*\n", title))

	if subtitle != "" {
		sb.WriteString(fmt.Sprintf("_%s_\n", subtitle))
	}

	sb.WriteString("\n")

	for i, plan := range plans {
		timeStr := "Весь день"
		if !plan.IsAllDay && plan.Time != "" {
			timeStr = "🕐 " + plan.Time
		}

		sb.WriteString(fmt.Sprintf("*%d.* %s\n", i+1, plan.Title))
		sb.WriteString(fmt.Sprintf("   _%s_\n", plan.Description))
		sb.WriteString(fmt.Sprintf("   %s\n\n", timeStr))
	}

	return sb.String()
}

func (b *Bot) createNavigationKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Сегодня", "today"),
			tgbotapi.NewInlineKeyboardButtonData("📅 Завтра", "tomorrow"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📆 Эта неделя", "week"),
			tgbotapi.NewInlineKeyboardButtonData("📆 След. неделя", "nextweek"),
		),
	)
}

func formatDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.Format("02.01.2006")
}
