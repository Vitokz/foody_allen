package commands

import (
	"context"
	"diet_bot/internal/entity"
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commands) MenuHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	messageText := "🍽️ *Главное меню* 🍽️\n\nВыберите действие из меню ниже:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥗 Создать рацион", flow.CommandGenerateDiet),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👀 Посмотреть рацион", flow.CommandSeeDiet),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛒 Посмотреть продукты", flow.CommandSeeDietProducts),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Заполнить конфигурацию", flow.CommandFillConfig),
		),
	)

	// Создаем сообщение с клавиатурой
	msg := tgbotapi.NewMessage(meta.ChatID, messageText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	return msg
}
