package creatediet

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
)

var foodCategoryFlow = []FillConfig{
	{
		// То что мы используем в рамках текущего состояния.
		State: flow.StateCreateFoodConfiguration_BreakfastCereals,
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.BreakfastCereals = value
			return nil
		},
		NextEvent: flow.EventCreateFoodConfigurationSideDishCereals,

		// То что мы используем если это состояние следующее.
		PromptText:  "Крупы для завтрака. Опишите какие крупы вы предпочитаете на завтрак",
		ExampleText: "овсянка,кукурузная каша",
	},
	{
		State:       flow.StateCreateFoodConfiguration_SideDishCereals,
		NextEvent:   flow.EventCreateFoodConfigurationVegetables,
		PromptText:  "Крупы для обеда. Опишите какие крупы вы предпочитаете на обед",
		ExampleText: "рис,гречка,булгур",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.SideDishCereals = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_Vegetables,
		NextEvent:   flow.EventCreateFoodConfigurationFruits,
		PromptText:  "Овощи. Опишите какие овощи вы предпочитаете",
		ExampleText: "капуста,брокколи,кабачки",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.Vegetables = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_Fruits,
		NextEvent:   flow.EventCreateFoodConfigurationNutsAndSeeds,
		PromptText:  "Фрукты. Опишите какие фрукты вы предпочитаете",
		ExampleText: "яблоко,груша,апельсин",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.Fruits = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_NutsAndSeeds,
		NextEvent:   flow.EventCreateFoodConfigurationDairyProducts,
		PromptText:  "Орехи и семечки. Опишите какие орехи и семечки вы предпочитаете",
		ExampleText: "миндаль,кешью,арахис",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.NutsAndSeeds = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_DairyProducts,
		NextEvent:   flow.EventCreateFoodConfigurationBread,
		PromptText:  "Молочные продукты. Опишите какие молочные продукты вы предпочитаете",
		ExampleText: "молоко,кефир,сыр",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.DairyProducts = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_Bread,
		NextEvent:   flow.EventCreateFoodConfigurationFish,
		PromptText:  "Хлеб. Опишите какой хлеб вы предпочитаете",
		ExampleText: "белый,черный,пшеничный",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.Bread = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_Fish,
		NextEvent:   flow.EventCreateFoodConfigurationMeat,
		PromptText:  "Рыба. Опишите какую рыбу вы предпочитаете",
		ExampleText: "лосось,треска,окунь",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.Fish = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_Meat,
		NextEvent:   flow.EventCreateFoodConfigurationEggs,
		PromptText:  "Мясо. Опишите какое мясо вы предпочитаете",
		ExampleText: "говядина,курица,телятина",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			config.FoodConfiguration.Meat = value
			return nil
		},
	},
	{
		State:       flow.StateCreateFoodConfiguration_Eggs,
		NextEvent:   flow.EventMainMenu,
		PromptText:  "Добавлять ли яйца в рацион? Ответьте \"да\" или \"нет\"",
		ExampleText: "да",
		FieldSetter: func(config *entity.DietConfiguration, value string) error {
			if value == "да" {
				config.FoodConfiguration.Eggs = true
			} else if value == "нет" {
				config.FoodConfiguration.Eggs = false
			} else {
				return fmt.Errorf("invalid value for eggs: %s", value)
			}
			return nil
		},
	},
}

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
