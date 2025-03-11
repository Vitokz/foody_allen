package creatediet

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

func (c *Commands) SideDishCerealsHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Крупы для обеда. Опишите какие крупы вы предпочитаете на обед

Например: рис,гречка,булгур
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationSideDishCereals); err != nil {
		c.logger.Errorw("Error creating food configuration side dish cereals", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) VegetablesHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Опишите какие овощи вы предпочитаете

Например: морковь,картофель,лук,помидоры,капуста,свекла,кабачок,тыква,баклажан
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationVegetables); err != nil {
		c.logger.Errorw("Error creating food configuration vegetables", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) FruitsHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Опишите какие фрукты вы предпочитаете

Например: яблоко,груша,банан
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationFruits); err != nil {
		c.logger.Errorw("Error creating food configuration fruits", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) NutsAndSeedsHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Опишите какие орехи и семечки вы предпочитаете

Например: грецкие орехи,семена тыквы,кешью,миндаль
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationNutsAndSeeds); err != nil {
		c.logger.Errorw("Error creating food configuration nuts and seeds", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) DairyProductsHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Опишите какие молочные продукты вы предпочитаете

Например: молоко,сыр,творог,греческиййогурт,кефир
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationDairyProducts); err != nil {
		c.logger.Errorw("Error creating food configuration dairy products", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) BreadHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Опишите какие хлебные изделия вы предпочитаете

Например: белый хлеб,бездрожжевой хлеб
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationBread); err != nil {
		c.logger.Errorw("Error creating food configuration bread", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) FishHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Опишите какие рыбы и морепродукты вы предпочитаете. Также сойдут ответы "не люблю рыбу" или "люблю рыбу, особенно язь"

Например: лосось,треска,креветки,кальмары,анчоусы. 
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationFish); err != nil {
		c.logger.Errorw("Error creating food configuration fish", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) MeatHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Опишите какие мясные продукты вы предпочитаете. Можно также ответить "не люблю мясо" или "люблю мясо, особенно говядину"

Например: говядина,курица
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationMeat); err != nil {
		c.logger.Errorw("Error creating food configuration meat", "error", err)

		return msg, err
	}

	return msg, nil
}

func (c *Commands) EggsHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Добавлять ли яйца в рацион? Ответьте "да" или "нет"

Например: да
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfigurationEggs); err != nil {
		c.logger.Errorw("Error creating food configuration eggs", "error", err)

		return msg, err
	}

	return msg, nil
}
