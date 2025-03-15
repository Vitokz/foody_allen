package creatediet

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

// Define a food category configuration struct
type FillConfig struct {
	State       string
	NextEvent   string
	PromptText  string
	ExampleText string
	FieldSetter func(config *entity.DietConfiguration, value string) error
}

var dietFillmentFlow = []FillConfig{
	{
		State:     "-",
		NextEvent: flow.EventCreateDiet,
		PromptText: `🚀 Ну штош начнем

Опиши свой образ жизни своими словами в свободной форме.`,
		ExampleText: `Работаю в офисе, сижу за компьютером. В день прохожу примерно 3 км. 
Занимаюсь кроссфитом 4 раза в неделю, по 2-3 часа каждая тренировка.`,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.LifestyleAndWorkouts = value
			return nil
		},
	},
	{
		State:     flow.StateCreateDiet_LifeStyle,
		NextEvent: flow.EventCreateDietTimeRestrictions,
		PromptText: `🚀 Ну штош начнем

Опиши свой образ жизни своими словами в свободной форме.`,
		ExampleText: `Работаю в офисе, сижу за компьютером. В день прохожу примерно 3 км. 
Занимаюсь кроссфитом 4 раза в неделю, по 2-3 часа каждая тренировка.`,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.LifestyleAndWorkouts = value
			return nil
		},
	},
	{
		State:       flow.StateCreateDiet_TimeRestrictions,
		NextEvent:   flow.EventCreateDietPFC,
		PromptText:  `Какие у тебя есть ограничения по времени?`,
		ExampleText: `Тренировка с 18:00 до 21:00, по-этому нет возможности есть в это время.`,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.TimeRestrictions = value
			return nil
		},
	},
	{
		State:     flow.StateCreateDiet_PFC,
		NextEvent: flow.EventCreateDietCalories,
		PromptText: `Введи количество белков, жиров и углеводов в процентах.

Формат: 30/30/40`,
		ExampleText: `20/30/50`,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			if err := config.PFC.ParsePFC(value); err != nil {
				return err
			}
			return nil
		},
	},
	{
		State:     flow.StateCreateDiet_Calories,
		NextEvent: flow.EventCreateDietNutritionPrinciples,
		PromptText: `Какое количество калорий тебе нужно в день?

Формат: 2000`,
		ExampleText: `3250`,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			calories, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			config.Calories = calories
			return nil
		},
	},
	{
		State:     flow.StateCreateDiet_NutritionPrinciples,
		NextEvent: flow.EventCreateDietIndividualRestrictions,
		PromptText: `Какие принципы питания ты хочешь поддерживать?

Формат свободный, но чтобы было понятно, что это за принципы.`,
		ExampleText: `1. Без перекусов.
2. Дробное питание на 5 приемов.
3. Разнообразие рациона.
4. Упор на медленные углеводы и овощи.
5. Завтрак - самый главный прием пищи.`,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.NutritionPrinciples = value
			return nil
		},
	},
	{
		State:     flow.StateCreateDiet_IndividualRestrictions,
		NextEvent: flow.EventCreateFoodConfiguration,
		PromptText: `Какие у тебя есть ограничения?

Формат свободный, но чтобы было понятно, что это за ограничения.`,
		ExampleText: `1. Не есть после 18:00.
2. Не есть сладкое.
3. Не использовать бобовые
4. Не использовать зеленый лук и болгарский перец.
5. Должно быть два вида гарнира, чтобы не приходилось готовить слишком много
6. Должен быть 1 тип рыбы. Чтобы не готовить слишком много
7. Должен быть 1 тип мяса. Чтобы не готовить слишком много
8. Два вида тушеных овощей вместо множества отдельных блюд. Чтобы не готовить слишком много.`,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.IndividualRestrictions = value
			return nil
		},
	},
}

// Map state to config for quick lookup
var flowConfigMap map[string]FillConfig

func init() {
	flowConfigMap = make(map[string]FillConfig)
	for _, config := range dietFillmentFlow {
		flowConfigMap[config.State] = config
	}
	for _, config := range foodCategoryFlow {
		flowConfigMap[config.State] = config
	}
}

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
	return c.executeFillment(ctx, flowConfigMap, update)
}

func (c *Commands) executeFillment(
	_ context.Context,
	flowConfig map[string]FillConfig,
	update *tgbotapi.Update,
) tgbotapi.Chattable {
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

	// Get the config for current state
	config, exists := flowConfig[chat.State]
	if !exists {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при формировании продуктов")
	}

	// Set the field value
	err = config.FieldSetter(dietConfiguration, meta.Message.Text)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, err.Error())
	}

	if config.NextEvent != flow.EventMainMenu {
		if err := botFSM.Event(config.NextEvent); err != nil {
			c.logger.Errorw("Error transitioning to next food category", "error", err)
			return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при переходе к следующему шагу")
		}
	} else {
		botFSM.SetState(flow.StateMenu)
	}

	// Save configuration
	err = c.saveDietConfiguration(chat, dietConfiguration)
	if err != nil {
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	// If this was the last category, complete the flow
	if config.NextEvent == flow.EventMainMenu {
		return tgbotapi.NewMessage(meta.ChatID, "Конфигурация завершена!")
	}

	// Otherwise, transition to the next category
	nextConfig := flowConfig[chat.State]
	// Return prompt for next category
	msg := makeResponseMsg(meta, nextConfig.PromptText, nextConfig.ExampleText)

	return msg
}

func makeResponseMsg(
	meta *entity.TelegramMeta,
	mainText string,
	exampleText string,
) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(meta.ChatID, "")
	msg.ParseMode = "MarkdownV2"

	mainText = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, mainText)
	exampleText = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, exampleText)

	msg.Text = fmt.Sprintf("\n❓ %s\n\nПример:\n```\n%s\n```\n", mainText, exampleText)

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
