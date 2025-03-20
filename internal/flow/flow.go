package flow

import (
	"context"

	"github.com/looplab/fsm"

	"diet_bot/internal/entity"
)

const (
	CommandStart        = "/start"
	CommandGenerateDiet = "/generate_diet"
	CommandMainMenu     = "/main_menu"
	CommandSeeDiet      = "/see_diet"
	CommandSeeDietDay   = "/see_diet_day_"
)

const (
	StateStart = "start"
	StateMenu  = "menu"
)

const (
	EventMainMenu = "main_menu"
)

const (
	StateCreateDiet                        = "create_diet"
	StateCreateDiet_LifeStyle              = "create_diet_life_style"
	StateCreateDiet_TimeRestrictions       = "create_diet_time_restrictions"
	StateCreateDiet_PFC                    = "create_diet_pfc"
	StateCreateDiet_Calories               = "create_diet_calories"
	StateCreateDiet_NutritionPrinciples    = "create_diet_nutrition_principles"
	StateCreateDiet_IndividualRestrictions = "create_diet_individual_restrictions"
	StateCreateDiet_FoodConfiguration      = "create_diet_food_configuration"
)

const (
	StateCreateFoodConfiguration_BreakfastCereals = "create_food_configuration_breakfast_cereals"
	StateCreateFoodConfiguration_SideDishCereals  = "create_food_configuration_side_dish_cereals"
	StateCreateFoodConfiguration_Vegetables       = "create_food_configuration_vegetables"
	StateCreateFoodConfiguration_Fruits           = "create_food_configuration_fruits"
	StateCreateFoodConfiguration_NutsAndSeeds     = "create_food_configuration_nuts_and_seeds"
	StateCreateFoodConfiguration_DairyProducts    = "create_food_configuration_dairy_products"
	StateCreateFoodConfiguration_Bread            = "create_food_configuration_bread"
	StateCreateFoodConfiguration_Fish             = "create_food_configuration_fish"
	StateCreateFoodConfiguration_Meat             = "create_food_configuration_meat"
	StateCreateFoodConfiguration_Eggs             = "create_food_configuration_eggs"
)

type BotFSM struct {
	FSM  *fsm.FSM
	Chat *entity.Chat
}

func NewBotFSM(chat *entity.Chat) *BotFSM {
	events := fsm.Events{
		{
			Name: EventMainMenu,
			Src:  []string{},
			Dst:  StateMenu,
		},
	}
	events = append(events, CreateDietFlow()...)
	events = append(events, CreateFoodConfigurationFlow()...)

	return &BotFSM{
		FSM: fsm.NewFSM(
			chat.State,
			events,
			fsm.Callbacks{},
		),
		Chat: chat,
	}
}

func (b *BotFSM) Event(event string) error {
	if err := b.FSM.Event(context.Background(), event, b.Chat); err != nil {
		return err
	}

	b.Chat.State = b.FSM.Current()

	return nil
}

func (b *BotFSM) SetState(event string) {
	b.FSM.SetState(event)
	b.Chat.State = b.FSM.Current()
}
