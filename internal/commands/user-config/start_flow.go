package userconfig

import (
	"context"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commands) StartFlow(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("error getting chat", "error", err)
		return nil, err
	}

	botFSM := flow.NewBotFSM(chat)
	if err := botFSM.Event(flow.EventUserConfigurationHeight); err != nil {
		c.logger.Errorw("error transitioning to start fill user config", "error", err)
		return nil, err
	}

	if err := c.repository.UpsertChat(chat); err != nil {
		c.logger.Errorw("error saving chat", "error", err)
		return nil, err
	}

	nextConfig := flowConfigMap[chat.State]
	return nextConfig.PromptMessage(meta.ChatID), nil
}
