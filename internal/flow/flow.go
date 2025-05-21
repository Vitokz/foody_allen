package flow

import (
	"context"

	"github.com/looplab/fsm"

	"diet_bot/internal/entity"
)

const (
	CommandStart               = "/start"
	CommandMenu                = "/menu"
	CommandGenerateDiet        = "/generate_diet"
	CommandGenerateDietDays    = "/generate_diet_days_"
	CommandSeeDiet             = "/see_diet"
	CommandSeeDietDay          = "/see_diet_day_"
	CommandSeeDietProducts     = "/see_diet_products"
	CommandFillConfig          = "/fill_config"
	CommandStartFillUserConfig = "/start_fill_user_config"

	CommandFillUserConfigGender       = "/fill_user_config_gender_"
	CommandFillUserConfigGenderMale   = "/fill_user_config_gender_male"
	CommandFillUserConfigGenderFemale = "/fill_user_config_gender_female"

	CommandFillDietGoal               = "/fill_diet_goal_"
	CommandFillDietGoalLoseWeight     = "/fill_diet_goal_lose_weight"
	CommandFillDietGoalMaintainWeight = "/fill_diet_goal_maintain_weight"
	CommandFillDietGoalGainWeight     = "/fill_diet_goal_gain_weight"

	CommandFillUserConfigActivity          = "/fill_user_config_activity_"
	CommandFillUserConfigActivitySedentary = "/fill_user_config_activity_sedentary"
	CommandFillUserConfigActivityLow       = "/fill_user_config_activity_low"
	CommandFillUserConfigActivityMedium    = "/fill_user_config_activity_medium"
	CommandFillUserConfigActivityHigh      = "/fill_user_config_activity_high"
	CommandFillUserConfigActivityVeryHigh  = "/fill_user_config_activity_very_high"

	CommandFillUserConfigDietType              = "/fill_user_config_diet_type_"
	CommandFillUserConfigDietTypeAnything      = "/fill_user_config_diet_type_anything"
	CommandFillUserConfigDietTypeKeto          = "/fill_user_config_diet_type_keto"
	CommandFillUserConfigDietTypePaleo         = "/fill_user_config_diet_type_paleo"
	CommandFillUserConfigDietTypeVegan         = "/fill_user_config_diet_type_vegan"
	CommandFillUserConfigDietTypeVegetarian    = "/fill_user_config_diet_type_vegetarian"
	CommandFillUserConfigDietTypeMediterranean = "/fill_user_config_diet_type_mediterranean"
)

const (
	StateStart = "start"
	StateMenu  = "menu"
)

const (
	EventMainMenu = "main_menu"
)

const (
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

const (
	StateUserConfiguration_Height    = "user_configuration_height"
	StateUserConfiguration_Weight    = "user_configuration_weight"
	StateUserConfiguration_Gender    = "user_configuration_gender"
	StateUserConfiguration_Age       = "user_configuration_age"
	StateUserConfiguration_Goal      = "user_configuration_goal"
	StateUserConfiguration_Activity  = "user_configuration_activity"
	StateUserConfiguration_DietType  = "user_configuration_diet_type"
	StateUserConfiguration_Allergies = "user_configuration_allergies"
	StateUserConfiguration_MealTypes = "user_configuration_meal_types"
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
	events = append(events, CreateUserConfigurationFlow()...)

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
	if event == EventMainMenu {
		b.FSM.SetState(StateMenu)
		b.Chat.State = StateMenu

		return nil
	}

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
