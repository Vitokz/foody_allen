package commands

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	generatediet "diet_bot/internal/commands/generate-diet"
	seediet "diet_bot/internal/commands/see-diet"
	userconfig "diet_bot/internal/commands/user-config"

	"diet_bot/internal/entity"
	internalerrors "diet_bot/internal/entity/errors"
	"diet_bot/internal/flow"
	"diet_bot/internal/repository"
)

type Commands struct {
	*generatediet.Command
	*seediet.DietCommand
	*userconfig.Commands
	logger     *zap.SugaredLogger
	repository *repository.Client
}

func NewCommands(
	repository *repository.Client,
	aiClient generatediet.AIClient,
	logger *zap.SugaredLogger,
) *Commands {
	return &Commands{
		Command:     generatediet.NewCommand(repository, aiClient, logger),
		DietCommand: seediet.NewDietCommand(repository, logger),
		Commands:    userconfig.NewCommands(repository, logger),
		repository:  repository,
		logger:      logger,
	}
}

const (
	StartCommand = `
👋 Привет! Я FoodyAlen - твой помощник по питанию

Я помогу тебе:
• Составить сбалансированный рацион на день
• Планировать питание на неделю
• Следить за БЖУ и калориями

Чтобы начать мне нужно будет немного узнать о тебе. Тыкай на кнопочку
`
)

func (c *Commands) StartHandler(ctx context.Context, update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	meta := entity.NewMeta(update)

	user, err := c.repository.GetUser(meta.Message.From.ID)
	if err != nil && err != internalerrors.ErrorUserNotFound {
		c.logger.Errorw("Error getting user", "error", err)

		return nil, err
	}

	if err == nil && user != nil {
		return c.MenuHandler(ctx, update)
	}

	user = &entity.User{
		ID:           meta.Message.From.ID,
		Username:     meta.Message.From.UserName,
		FirstName:    meta.Message.From.FirstName,
		LastName:     meta.Message.From.LastName,
		LanguageCode: meta.Message.From.LanguageCode,
		IsBot:        meta.Message.From.IsBot,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	chat := &entity.Chat{
		ID:     meta.ChatID,
		State:  flow.StateStart,
		UserID: user.ID,
	}

	userConfig := &entity.UserConfiguration{
		ID:     uuid.New(),
		UserID: user.ID,
	}

	err = c.repository.UpsertUser(user)
	if err != nil {
		c.logger.Errorw("Error saving user", "error", err)

		return nil, err
	}

	err = c.repository.UpsertChat(chat)
	if err != nil {
		c.logger.Errorw("Error saving chat", "error", err)

		return nil, err
	}

	err = c.repository.CreateUserConfiguration(userConfig)
	if err != nil {
		c.logger.Errorw("Error saving user configuration", "error", err)

		return nil, err
	}

	msg := tgbotapi.NewMessage(meta.ChatID, StartCommand)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Заполнить информацию", flow.CommandStartFillUserConfig),
		),
	)

	return msg, nil
}
