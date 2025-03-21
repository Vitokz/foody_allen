package generatediet

import (
	"context"
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

type AIClient interface {
	GenerateDiet(prompt string) (string, error)
}

type Repository interface {
	GetDietConfiguration(userID int64) (*entity.DietConfiguration, error)
	CreateDiet(diet *entity.GeneratedDiet) error
}

type Command struct {
	repository Repository
	aiClient   AIClient
	logger     *zap.SugaredLogger
}

func NewCommand(repository Repository, aiClient AIClient, logger *zap.SugaredLogger) *Command {
	return &Command{
		repository: repository,
		aiClient:   aiClient,
		logger:     logger,
	}
}

func (c *Command) GenerateDietHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	configuration, err := c.repository.GetDietConfiguration(meta.UserID)
	if err != nil {
		c.logger.Error("error getting diet configuration", zap.Error(err))
		return nil
	}

	prompt := GenerateDietPrompt(configuration)

	c.logger.Info("generated diet prompt", zap.String("prompt", prompt))

	response, err := c.aiClient.GenerateDiet(prompt)
	if err != nil {
		c.logger.Error("error generating diet", zap.Error(err))
		return nil
	}

	c.logger.Info("generated diet", zap.String("diet", response))

	diet := new(entity.GeneratedDiet)
	err = json.Unmarshal([]byte(response), diet)
	if err != nil {
		c.logger.Error("error unmarshalling diet", zap.Error(err))
		return nil
	}

	diet.SetIDs()
	diet.UserID = meta.UserID

	err = c.repository.CreateDiet(diet)
	if err != nil {
		c.logger.Error("error creating diet", zap.Error(err))
		return nil
	}

	msg := tgbotapi.NewMessage(meta.ChatID, "Рацион успешно создан!")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥗 Посмотреть рацион", flow.CommandSeeDiet),
		),
	)

	return msg
}
