package userconfig

import (
	"context"
	"diet_bot/internal/entity"
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commands) ProfileMenu(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("error getting chat", "error", err)
		return nil, err
	}

	userConfig, err := c.repository.GetUserConfiguration(chat.UserID)
	if err != nil {
		c.logger.Errorw("error getting user configuration", "error", err)
		return nil, err
	}

	msg := tgbotapi.NewMessage(meta.ChatID, userConfig.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Заполнить заново", flow.CommandStartFillUserConfig),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", flow.CommandMenu),
		),
	)

	return msg, nil
}
