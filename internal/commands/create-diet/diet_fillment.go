package creatediet

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

func (c *Commands) IsFillDiet(ctx context.Context, update *tgbotapi.Update) bool {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return false
	}

	ok := flow.IsCreateDietFillState(chat.State)

	if !ok {
		c.logger.Infof("Chat state: %s is not fill diet", chat.State)
	} else {
		c.logger.Infof("Chat state: %s is fill diet", chat.State)
	}

	return ok
}

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
		return c.timeRestrictionsFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateDiet_PFC:
		return c.pfcFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateDiet_Calories:
		return c.caloriesFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateDiet_NutritionPrinciples:
		return c.nutritionPrinciplesFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateDiet_IndividualRestrictions:
		return c.individualRestrictionsFillment(ctx, botFSM, meta, chat, dietConfiguration)
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

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) timeRestrictionsFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	dietConfiguration.TimeRestrictions = meta.Message.Text

	msg, err := c.PFCHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		c.logger.Errorw("Error creating diet pfc", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при создании диеты")
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) pfcFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	if err := dietConfiguration.PFC.ParsePFC(meta.Message.Text); err != nil {
		c.logger.Errorw("Error parsing PFC", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Неправильный формат БЖУ")
	}

	msg, err := c.CaloriesHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		c.logger.Errorw("Error creating diet calories", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при создании диеты")
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) caloriesFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	calories, err := strconv.Atoi(meta.Message.Text)
	if err != nil {
		c.logger.Errorw("Error converting calories to int", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Неправильный формат количества калорий")
	}

	dietConfiguration.Calories = calories

	msg, err := c.NutritionPrinciplesHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		c.logger.Errorw("Error creating diet nutrition principles", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при создании диеты")
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) nutritionPrinciplesFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	dietConfiguration.NutritionPrinciples = meta.Message.Text

	msg, err := c.IndividualRestrictionsHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		c.logger.Errorw("Error creating diet individual restrictions", "error", err)
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		c.logger.Errorw("Error saving diet configuration", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) individualRestrictionsFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	dietConfiguration.IndividualRestrictions = meta.Message.Text

	msg, err := c.FoodConfigurationHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		c.logger.Errorw("Error creating diet food configuration", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при создании диеты")
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) saveDietConfiguration(
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) error {
	err := c.repository.UpsertDietConfiguration(dietConfiguration)
	if err != nil {
		c.logger.Errorw("Error upserting diet configuration", "error", err)
		return err
	}

	err = c.repository.UpsertChat(chat)
	if err != nil {
		c.logger.Errorw("Error upserting chat", "error", err)
		return err
	}

	return nil
}
