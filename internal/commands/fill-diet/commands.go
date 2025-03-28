package creatediet

import (
	"context"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	GetChat(id int64) (*entity.Chat, error)
	UpsertChat(chat *entity.Chat) error

	GetDietConfiguration(userID int64) (*entity.DietConfiguration, error)
	UpsertDietConfiguration(dietConfiguration *entity.DietConfiguration) error
}

type Commands struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func NewCommands(repository Repository, logger *zap.SugaredLogger) *Commands {
	return &Commands{
		repository: repository,
		logger:     logger,
	}
}

func (c *Commands) CreateDietHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	startFlowStep := dietFillmentFlow[0]

	msg := makeResponseMsg(meta, startFlowStep.PromptText, startFlowStep.ExampleText)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	dietConfiguration, err := c.repository.GetDietConfiguration(chat.UserID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			c.logger.Errorw("Error getting diet configuration", "error", err)
			return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
		}

		dietConfiguration = &entity.DietConfiguration{
			UserID: chat.UserID,
		}
	}

	botFSM := flow.NewBotFSM(chat)

	botFSM.Event(startFlowStep.NextEvent)

	// Save configuration
	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}
