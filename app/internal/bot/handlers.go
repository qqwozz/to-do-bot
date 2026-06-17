package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	state := b.convs.get(chatID)

	switch state.step {
	case stepTitle:
		state.plan.Title = strings.TrimSpace(msg.Text)
		state.step = stepDesc
		b.convs.set(chatID, state)
		b.sendWithKb(chatID, "📝 Описание:", cancelKeyboard())

	case stepDesc:
		state.plan.Description = strings.TrimSpace(msg.Text)
		state.step = stepDate
		b.convs.set(chatID, state)
		b.sendWithKb(chatID, "🗓 Дата:", createDateKeyboard())

	case stepDate:
		t, err := time.Parse("2006-01-02", strings.TrimSpace(msg.Text))
		if err != nil {
			b.sendWithKb(chatID, "❌ Формат: ГГГГ-ММ-ДД. Попробуйте ещё раз:", cancelKeyboard())
			return
		}
		state.plan.Date = t.Format("2006-01-02")
		state.plan.IsAllDay = true
		state.step = stepTime
		b.convs.set(chatID, state)
		b.sendWithKb(chatID, "⏰ Время:", createTimeKeyboard())

	case stepTime:
		t, err := time.Parse("15:04", strings.TrimSpace(msg.Text))
		if err != nil {
			b.sendWithKb(chatID, "❌ Формат: ЧЧ:ММ. Попробуйте ещё раз:", cancelKeyboard())
			return
		}
		state.plan.Time = t.Format("15:04")
		state.plan.IsAllDay = false
		b.convs.reset(chatID)
		b.createAndConfirm(chatID, state.plan)

	default:
		b.sendWithKb(chatID, "Используйте кнопки для навигации.", mainMenuKeyboard())
	}
}

func (b *Bot) handleCallback(q *tgbotapi.CallbackQuery) {
	chatID := q.Message.Chat.ID
	data := q.Data
	b.answerCallback(q.ID)

	// сброс диалога при любом action-кнопке
	if data != "plan:cancel" && strings.HasPrefix("plan:", data) && data != "plan:create" {
		if s := b.convs.get(chatID); s.step != stepIdle {
			b.convs.reset(chatID)
		}
	}

	switch {
	case data == "menu:main":
		b.sendMainMenu(chatID)

	case data == "plans:show":
		b.sendWithKb(chatID, "📅 Выберите период:", plansPeriodKeyboard())

	case data == "plans:today":
		b.sendPlansForDate(chatID, time.Now().Format("2006-01-02"), "Сегодня")

	case data == "plans:tomorrow":
		d := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		b.sendPlansForDate(chatID, d, "Завтра")

	case data == "plans:week":
		b.sendPlansForWeek(chatID, 0, "Эта неделя")

	case data == "plans:nextweek":
		b.sendPlansForWeek(chatID, 7, "Следующая неделя")

	case data == "plan:create":
		b.startPlanCreation(chatID)

	case data == "plan:cancel":
		b.convs.reset(chatID)
		b.sendMainMenu(chatID)

	case data == "plan:date:today":
		b.setDateAndAskTime(chatID, time.Now().Format("2006-01-02"))

	case data == "plan:date:tomorrow":
		d := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		b.setDateAndAskTime(chatID, d)

	case data == "plan:date:week":
		d := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
		b.setDateAndAskTime(chatID, d)

	case data == "plan:date:custom":
		b.sendWithKb(chatID, "✏️ Введите дату (ГГГГ-ММ-ДД):", cancelKeyboard())

	case strings.HasPrefix(data, "plan:time:"):
		b.setTimeAndCreate(chatID, data)

	case data == "menu:plans":
		b.sendWithKb(chatID, "📅 Выберите период:", plansPeriodKeyboard())
	}
}

func (b *Bot) startPlanCreation(chatID int64) {
	state := &conversation{step: stepTitle}
	b.convs.set(chatID, state)
	b.sendWithKb(chatID, "📝 Название:", cancelKeyboard())
}

func (b *Bot) setDateAndAskTime(chatID int64, date string) {
	state := b.convs.get(chatID)
	state.plan.Date = date
	state.plan.IsAllDay = true
	state.step = stepTime
	b.convs.set(chatID, state)
	b.sendWithKb(chatID, "⏰ Время:", createTimeKeyboard())
}

func (b *Bot) setTimeAndCreate(chatID int64, data string) {
	state := b.convs.get(chatID)
	switch data {
	case "plan:time:morning":
		state.plan.Time = "09:00"
	case "plan:time:afternoon":
		state.plan.Time = "14:00"
	case "plan:time:evening":
		state.plan.Time = "18:00"
	case "plan:time:allday":
		state.plan.IsAllDay = true
	}
	b.convs.reset(chatID)
	b.createAndConfirm(chatID, state.plan)
}

func (b *Bot) createAndConfirm(chatID int64, plan PlanRequest) {
	if err := b.client.CreatePlan(plan); err != nil {
		b.send(chatID, fmt.Sprintf("❌ Ошибка: %v", err))
		return
	}

	var sb strings.Builder
	sb.WriteString("✅ План создан!\n\n")
	sb.WriteString(fmt.Sprintf("📝 *%s*\n", plan.Title))
	sb.WriteString(fmt.Sprintf("   _%s_\n", plan.Description))
	sb.WriteString(fmt.Sprintf("   📅 %s\n", formatDate(plan.Date)))
	if !plan.IsAllDay && plan.Time != "" {
		sb.WriteString(fmt.Sprintf("   🕐 %s\n", plan.Time))
	} else {
		sb.WriteString("   ⏰ Весь день\n")
	}

	b.sendWithKb(chatID, sb.String(), mainMenuKeyboard())
}

func (b *Bot) sendMainMenu(chatID int64) {
	b.sendWithKb(chatID, "👋 *To-Do бот*\n\nВыберите действие:", mainMenuKeyboard())
}

func (b *Bot) sendPlansForDate(chatID int64, date, label string) {
	plans, err := b.client.GetPlansByDate(date)
	if err != nil {
		b.send(chatID, "❌ Ошибка получения планов")
		return
	}
	b.sendWithKb(chatID, formatPlansList(plans, label, ""), plansViewKeyboard())
}

func (b *Bot) sendPlansForWeek(chatID int64, offsetDays int, label string) {
	now := time.Now()
	start := now.AddDate(0, 0, -int(now.Weekday())+offsetDays)
	end := start.AddDate(0, 0, 6)

	plans, err := b.client.GetPlansByDateRange(start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		b.send(chatID, "❌ Ошибка получения планов")
		return
	}

	subtitle := fmt.Sprintf("%s — %s", formatDate(start.Format("2006-01-02")), formatDate(end.Format("2006-01-02")))
	b.sendWithKb(chatID, formatPlansList(plans, label, subtitle), plansViewKeyboard())
}
