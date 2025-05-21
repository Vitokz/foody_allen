package entity

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMeta struct {
	Message      *tgbotapi.Message
	CallbackData *string
	Chat         *tgbotapi.Chat
	ChatID       int64
	UserID       int64
}

func NewMeta(update *tgbotapi.Update) *TelegramMeta {
	msg := update.Message
	var callbackData *string

	isCallback := false
	if msg == nil {
		msg = update.CallbackQuery.Message
		callbackData = &update.CallbackQuery.Data
		isCallback = true
	}

	chatID := msg.Chat.ID
	userID := msg.From.ID
	if isCallback {
		userID = msg.Chat.ID
	}

	return &TelegramMeta{
		Message:      msg,
		CallbackData: callbackData,
		Chat:         msg.Chat,
		ChatID:       chatID,
		UserID:       userID,
	}
}
