package creatediet

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

func (c *Commands) FillDiet(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	dietConfiguration, err := c.repository.GetDietConfiguration(chat.UserID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			c.logger.Errorw("Error upserting diet configuration", "error", err)
			return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
		}

		dietConfiguration = &entity.DietConfiguration{
			ID:     uuid.New().String(),
			UserID: chat.UserID,
		}
	}

	botFSM := flow.NewBotFSM(chat)

	switch chat.State {
	case flow.StateCreateDiet_LifeStyle:
		return c.lifeStyleFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateDiet_TimeRestrictions:
	case flow.StateCreateDiet_PFC:
	case flow.StateCreateDiet_Calories:
	case flow.StateCreateDiet_NutritionPrinciples:
	case flow.StateCreateDiet_IndividualRestrictions:
	default:
	}

	return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при создании диеты")
}

func (c *Commands) lifeStyleFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	dietConfiguration.LifestyleAndWorkouts = meta.Message.Text

	msg, err := c.TimeRestrictionsHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		c.logger.Errorw("Error creating diet time restrictions", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при создании диеты")
	}

	err = c.repository.UpsertDietConfiguration(dietConfiguration)
	if err != nil {
		c.logger.Errorw("Error upserting diet configuration", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	err = c.repository.UpsertChat(chat)
	if err != nil {
		c.logger.Errorw("Error upserting chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}
