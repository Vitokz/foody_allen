package flow

import "github.com/looplab/fsm"

const (
	EventCreateFoodConfiguration                 = "create_food_configuration"
	EventCreateFoodConfigurationBreakfastCereals = "create_food_configuration_breakfast_cereals"
	EventCreateFoodConfigurationSideDishCereals  = "create_food_configuration_side_dish_cereals"
	EventCreateFoodConfigurationVegetables       = "create_food_configuration_vegetables"
	EventCreateFoodConfigurationFruits           = "create_food_configuration_fruits"
	EventCreateFoodConfigurationNutsAndSeeds     = "create_food_configuration_nuts_and_seeds"
	EventCreateFoodConfigurationDairyProducts    = "create_food_configuration_dairy_products"
	EventCreateFoodConfigurationBread            = "create_food_configuration_bread"
	EventCreateFoodConfigurationFish             = "create_food_configuration_fish"
	EventCreateFoodConfigurationMeat             = "create_food_configuration_meat"
	EventCreateFoodConfigurationEggs             = "create_food_configuration_eggs"
)

func CreateFoodConfigurationFlow() fsm.Events {
	return fsm.Events{
		{
			Name: EventCreateFoodConfiguration,
			Src:  []string{StateMenu, StateCreateDiet_IndividualRestrictions},
			Dst:  StateCreateFoodConfiguration_BreakfastCereals,
		},
		{
			Name: EventCreateFoodConfigurationSideDishCereals,
			Src:  []string{StateCreateFoodConfiguration_BreakfastCereals},
			Dst:  StateCreateFoodConfiguration_SideDishCereals,
		},
		{
			Name: EventCreateFoodConfigurationVegetables,
			Src:  []string{StateCreateFoodConfiguration_SideDishCereals},
			Dst:  StateCreateFoodConfiguration_Vegetables,
		},
		{
			Name: EventCreateFoodConfigurationFruits,
			Src:  []string{StateCreateFoodConfiguration_Vegetables},
			Dst:  StateCreateFoodConfiguration_Fruits,
		},
		{
			Name: EventCreateFoodConfigurationNutsAndSeeds,
			Src:  []string{StateCreateFoodConfiguration_Fruits},
			Dst:  StateCreateFoodConfiguration_NutsAndSeeds,
		},
		{
			Name: EventCreateFoodConfigurationDairyProducts,
			Src:  []string{StateCreateFoodConfiguration_NutsAndSeeds},
			Dst:  StateCreateFoodConfiguration_DairyProducts,
		},
		{
			Name: EventCreateFoodConfigurationBread,
			Src:  []string{StateCreateFoodConfiguration_DairyProducts},
			Dst:  StateCreateFoodConfiguration_Bread,
		},
		{
			Name: EventCreateFoodConfigurationFish,
			Src:  []string{StateCreateFoodConfiguration_Bread},
			Dst:  StateCreateFoodConfiguration_Fish,
		},
		{
			Name: EventCreateFoodConfigurationMeat,
			Src:  []string{StateCreateFoodConfiguration_Fish},
			Dst:  StateCreateFoodConfiguration_Meat,
		},
		{
			Name: EventCreateFoodConfigurationEggs,
			Src:  []string{StateCreateFoodConfiguration_Meat},
			Dst:  StateCreateFoodConfiguration_Eggs,
		},
	}
}

func IsCreateFoodConfigurationFillState(state string) bool {
	return state == StateCreateFoodConfiguration_BreakfastCereals ||
		state == StateCreateFoodConfiguration_SideDishCereals ||
		state == StateCreateFoodConfiguration_Vegetables ||
		state == StateCreateFoodConfiguration_Fruits ||
		state == StateCreateFoodConfiguration_NutsAndSeeds ||
		state == StateCreateFoodConfiguration_DairyProducts ||
		state == StateCreateFoodConfiguration_Fish ||
		state == StateCreateFoodConfiguration_Meat ||
		state == StateCreateFoodConfiguration_Eggs ||
		state == StateCreateFoodConfiguration_Bread
}
