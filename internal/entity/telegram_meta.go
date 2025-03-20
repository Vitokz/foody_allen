package entity

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMeta struct {
	Message *tgbotapi.Message
	Chat    *tgbotapi.Chat
	ChatID  int64
	UserID  int64
}

func NewMeta(update *tgbotapi.Update) *TelegramMeta {
	msg := update.Message

	if msg == nil {
		msg = update.CallbackQuery.Message
	}

	chatID := msg.Chat.ID
	userID := msg.From.ID

	fmt.Println(msg.Chat)
	fmt.Println(msg.From)
	fmt.Println("Смотри какой ID", msg.From.ID)
	fmt.Println("Смотри какой у чата ID", msg.Chat.ID)

	return &TelegramMeta{
		Message: msg,
		Chat:    msg.Chat,
		ChatID:  chatID,
		UserID:  userID,
	}
}
