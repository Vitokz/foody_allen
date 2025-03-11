package creatediet

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

func (c *Commands) IsFillProducts(ctx context.Context, update *tgbotapi.Update) bool {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return false
	}

	ok := flow.IsCreateFoodConfigurationFillState(chat.State)

	if !ok {
		c.logger.Infof("Chat state: %s is not fill products", chat.State)
	} else {
		c.logger.Infof("Chat state: %s is fill products", chat.State)
	}

	return ok
}

func (c *Commands) FillProducts(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	dietConfiguration, err := c.repository.GetDietConfiguration(chat.UserID)
	if err != nil {
		c.logger.Errorw("Error getting diet configuration", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	botFSM := flow.NewBotFSM(chat)

	switch chat.State {
	case flow.StateCreateFoodConfiguration_BreakfastCereals:
		return c.breakfastCerealsFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_SideDishCereals:
		return c.sideDishCerealsFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_Vegetables:
		return c.vegetablesFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_Fruits:
		return c.fruitsFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_NutsAndSeeds:
		return c.nutsAndSeedsFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_DairyProducts:
		return c.dairyProductsFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_Bread:
		return c.breadFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_Fish:
		return c.fishFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_Meat:
		return c.meatFillment(ctx, botFSM, meta, chat, dietConfiguration)
	case flow.StateCreateFoodConfiguration_Eggs:
		return c.eggsFillment(ctx, botFSM, meta, chat, dietConfiguration)
	default:
	}

	return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при формировании продуктов")
}

func (c *Commands) breakfastCerealsFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.BreakfastCereals = meta.Message.Text

	msg, err := c.SideDishCerealsHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) sideDishCerealsFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.SideDishCereals = meta.Message.Text

	msg, err := c.VegetablesHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) vegetablesFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.Vegetables = meta.Message.Text

	msg, err := c.FruitsHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) fruitsFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.Fruits = meta.Message.Text

	msg, err := c.NutsAndSeedsHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) nutsAndSeedsFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.NutsAndSeeds = meta.Message.Text

	msg, err := c.DairyProductsHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) dairyProductsFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.DairyProducts = meta.Message.Text

	msg, err := c.BreadHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) breadFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.Bread = meta.Message.Text

	msg, err := c.FishHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) fishFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.Fish = meta.Message.Text

	msg, err := c.MeatHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) meatFillment(ctx context.Context, botFSM *flow.BotFSM, meta *entity.TelegramMeta, chat *entity.Chat, dietConfiguration *entity.DietConfiguration) tgbotapi.Chattable {
	dietConfiguration.FoodConfiguration.Meat = meta.Message.Text

	msg, err := c.EggsHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) eggsFillment(
	ctx context.Context,
	botFSM *flow.BotFSM,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
) tgbotapi.Chattable {
	if meta.Message.Text == "да" {
		dietConfiguration.FoodConfiguration.Eggs = true
	} else if meta.Message.Text == "нет" {
		dietConfiguration.FoodConfiguration.Eggs = false
	} else {
		return tgbotapi.NewMessage(meta.ChatID, "Пожалуйста, ответьте 'да' или 'нет'")
	}

	msg, err := c.CompleteConfigurationHandler(ctx, meta, chat, dietConfiguration, botFSM)
	if err != nil {
		return msg
	}

	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}
