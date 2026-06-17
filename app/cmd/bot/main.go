package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"todo-bot/internal/bot"
	"todo-bot/internal/config"
)

func main() {
	log.Println("[BOOT] running tests...")
	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("[BOOT] tests failed: %v", err)
	}
	log.Println("[BOOT] tests passed")

	cfg := config.Load()

	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	tgBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("create bot: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go tgBot.Run()

	<-sig
	log.Println("shutting down")
}
