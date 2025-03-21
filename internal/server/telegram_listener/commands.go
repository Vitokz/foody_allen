package telegramlistener

import (
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setupBotCommands(bot *tgbotapi.BotAPI) error {
	commands := []tgbotapi.BotCommand{
		{Command: flow.CommandStart, Description: "Начать работу с ботом"},
		{Command: flow.CommandMenu, Description: "Главное меню"},
	}

	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(config)
	return err
}
