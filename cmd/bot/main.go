package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"todo-bot/internal/bot"
	"todo-bot/internal/config"
)

func main() {
	cfg := config.Load()

	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	tgBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go tgBot.Run()

	sig := <-sigChan
	log.Printf("Signal %v received, shutting down...", sig)
}
