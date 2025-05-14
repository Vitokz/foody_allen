package commands

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
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
👋 Привет! Я FoodyAlen - твой AI-помощник по питанию

Я помогу тебе:
• Составить сбалансированный рацион на день
• Планировать питание на неделю
• Следить за БЖУ и калориями

У тебя уже есть базовая конфигурация рациона, можешь сразу попробовать сформировать рацион! 
`
)

func (c *Commands) StartHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	user, err := c.repository.GetUser(meta.Message.From.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		c.logger.Errorw("Error getting user", "error", err)
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

	defaultDietConfiguration := entity.DefaultDietConfiguration()
	defaultDietConfiguration.UserID = user.ID

	err = c.repository.UpsertUser(user)
	if err != nil {
		c.logger.Errorw("Error saving user", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	err = c.repository.UpsertChat(chat)
	if err != nil {
		c.logger.Errorw("Error saving chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	err = c.repository.UpsertDietConfiguration(defaultDietConfiguration)
	if err != nil {
		c.logger.Errorw("Error saving diet configuration", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	msg := tgbotapi.NewMessage(meta.ChatID, StartCommand)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥗 Создать рацион", flow.CommandGenerateDiet),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠  В меню", flow.CommandMenu),
		),
	)

	return msg
}
