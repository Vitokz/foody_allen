package commands

import (
	"context"
	"diet_bot/internal/entity"
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commands) MenuHandler(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
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
			tgbotapi.NewInlineKeyboardButtonData("👤 Профиль", flow.CommandProfileMenu),
		),
	)

	msg := tgbotapi.NewMessage(meta.ChatID, messageText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("error getting chat", "error", err)
		return nil, err
	}

	botFSM := flow.NewBotFSM(chat)
	if err := botFSM.Event(flow.EventMainMenu); err != nil {
		c.logger.Errorw("error transitioning to main menu", "error", err)
		return nil, err
	}

	if err := c.repository.UpsertChat(chat); err != nil {
		c.logger.Errorw("error saving chat", "error", err)
		return nil, err
	}

	return msg, nil
}
