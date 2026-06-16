package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"todo-bot/internal/config"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	client *Client
	config *config.Config
}

func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	log.Printf("[BOT] Авторизован как @%s", api.Self.UserName)

	return &Bot{
		api:    api,
		client: NewClient(cfg.BackendURL),
		config: cfg,
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	log.Printf("[BOT] Бот запущен...")

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
			continue
		}
		if update.CallbackQuery != nil {
			b.handleCallbackQuery(update.CallbackQuery)
		}
	}
}

func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("[BOT] Ошибка отправки: %v", err)
	}
}

func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("[BOT] Ошибка отправки: %v", err)
	}
}

func (b *Bot) AnswerCallback(queryID string) {
	callback := tgbotapi.NewCallback(queryID, "")
	b.api.Request(callback)
}
