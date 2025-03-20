package main

import (
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	"diet_bot/internal/clients/gpt"
	"diet_bot/internal/commands"
	"diet_bot/internal/lib/logger"
	"diet_bot/internal/repository"
	telegramlistener "diet_bot/internal/server/telegram_listener"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	logger, err := logger.NewLogger("development")
	if err != nil {
		panic(err)
	}

	database, err := repository.NewClient()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	aiClient := gpt.NewClient(os.Getenv("OPENAI_API_KEY"), os.Getenv("OPENAI_BASE_URL"))

	commands := commands.NewCommands(database, aiClient, logger)

	listener, err := telegramlistener.NewListener(os.Getenv("TELEGRAM_BOT_TOKEN"), commands, logger)
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go listener.Listen()

	<-quit

	logger.Info("Shutting down gracefully...")

	listener.Stop()
}
