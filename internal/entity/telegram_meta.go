package entity

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type TelegramMeta struct {
	Message *tgbotapi.Message
	Chat    *tgbotapi.Chat
	ChatID  int64
}

func NewMeta(update *tgbotapi.Update) *TelegramMeta {
	msg := update.Message

	if msg == nil {
		msg = update.CallbackQuery.Message
	}

	chatID := msg.Chat.ID

	return &TelegramMeta{
		Message: msg,
		Chat:    msg.Chat,
		ChatID:  chatID,
	}
}
