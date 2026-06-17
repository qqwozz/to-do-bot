package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"todo-bot/internal/config"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	client *Client
	convs  *conversationStore
}

func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	log.Printf("[BOT] authorized as @%s", api.Self.UserName)

	return &Bot{
		api:    api,
		client: NewClient(cfg.BackendURL),
		convs:  newConversationStore(),
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	log.Printf("[BOT] started")

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
			continue
		}
		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("[BOT] send error: %v", err)
	}
}

func (b *Bot) sendWithKb(chatID int64, text string, kb tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = kb
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("[BOT] send error: %v", err)
	}
}

func (b *Bot) answerCallback(id string) {
	b.api.Request(tgbotapi.NewCallback(id, ""))
}
