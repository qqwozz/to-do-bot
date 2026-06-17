package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func mainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать план", "plan:create"),
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои планы", "plans:show"),
		),
	)
}

func plansPeriodKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Сегодня", "plans:today"),
			tgbotapi.NewInlineKeyboardButtonData("📅 Завтра", "plans:tomorrow"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📆 Эта неделя", "plans:week"),
			tgbotapi.NewInlineKeyboardButtonData("📆 След. неделя", "plans:nextweek"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "menu:main"),
		),
	)
}

func plansViewKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 Сегодня", "plans:today"),
			tgbotapi.NewInlineKeyboardButtonData("📅 Завтра", "plans:tomorrow"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📆 Эта неделя", "plans:week"),
			tgbotapi.NewInlineKeyboardButtonData("📆 След. неделя", "plans:nextweek"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Создать план", "plan:create"),
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", "menu:main"),
		),
	)
}

func createDateKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗓 Сегодня", "plan:date:today"),
			tgbotapi.NewInlineKeyboardButtonData("🗓 Завтра", "plan:date:tomorrow"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗓 Через неделю", "plan:date:week"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Ввести дату", "plan:date:custom"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "plan:cancel"),
		),
	)
}

func createTimeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌅 Утро (09:00)", "plan:time:morning"),
			tgbotapi.NewInlineKeyboardButtonData("☀️ День (14:00)", "plan:time:afternoon"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌆 Вечер (18:00)", "plan:time:evening"),
			tgbotapi.NewInlineKeyboardButtonData("⏰ Весь день", "plan:time:allday"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "plan:cancel"),
		),
	)
}

func cancelKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "plan:cancel"),
		),
	)
}
