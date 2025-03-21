package commands

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	filldiet "diet_bot/internal/commands/fill-diet"
	generatediet "diet_bot/internal/commands/generate-diet"
	seediet "diet_bot/internal/commands/see-diet"
	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

type Commands struct {
	*filldiet.Commands
	*generatediet.Command
	*seediet.DietCommand
	logger     *zap.SugaredLogger
	repository Repository
}

func NewCommands(repository Repository, aiClient generatediet.AIClient, logger *zap.SugaredLogger) *Commands {
	return &Commands{
		Commands:    filldiet.NewCommands(repository, logger),
		Command:     generatediet.NewCommand(repository, aiClient, logger),
		DietCommand: seediet.NewDietCommand(repository, logger),
		repository:  repository,
		logger:      logger,
	}
}

const (
	StartCommand = `
Привет! Меня зовут FoodyAlen. Я личный помошник в составлении рациона.

Вот какие проблемы я могу решить:
1. Составить рацион на день с учетом БЖУ и Калорий. 
2. Составить рацион на неделю, чтобы в воскресенье все приготовить, упаковать, а потом съесть.

Чтобы помочь тебе в формировании рациона, мне нужно узнать тебя поближе 👉👈.

Давай заполним конфигурацию рациона, это не займет больше 5 минут.
`
)

func (c *Commands) StartHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	user := entity.User{
		ID:           meta.Message.From.ID,
		Username:     meta.Message.From.UserName,
		FirstName:    meta.Message.From.FirstName,
		LastName:     meta.Message.From.LastName,
		LanguageCode: meta.Message.From.LanguageCode,
		IsBot:        meta.Message.From.IsBot,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	chat := entity.Chat{
		ID:     meta.ChatID,
		State:  flow.StateStart,
		UserID: user.ID,
	}

	err := c.repository.UpsertUser(&user)
	if err != nil {
		c.logger.Errorw("Error saving user", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	err = c.repository.UpsertChat(&chat)
	if err != nil {
		c.logger.Errorw("Error saving chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	button := tgbotapi.NewInlineKeyboardButtonData("Заполнить конфигурацию", flow.CommandFillConfig)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
	msg := tgbotapi.NewMessage(meta.ChatID, StartCommand)
	msg.ReplyMarkup = keyboard

	return msg
}
