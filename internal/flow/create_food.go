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
	EventCreateFoodConfigurationFish             = "create_food_configuration_fish"
	EventCreateFoodConfigurationMeat             = "create_food_configuration_meat"
	EventCreateFoodConfigurationEggs             = "create_food_configuration_eggs"
)

func CreateFoodConfigurationFlow() fsm.Events {
	return fsm.Events{
		{
			Name: EventCreateFoodConfiguration,
			Src:  []string{StateMenu, StateCreateDiet_FoodConfiguration},
			Dst:  StateCreateFoodConfiguration,
		},
		{
			Name: EventCreateFoodConfigurationBreakfastCereals,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_BreakfastCereals,
		},
		{
			Name: EventCreateFoodConfigurationSideDishCereals,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_SideDishCereals,
		},
		{
			Name: EventCreateFoodConfigurationVegetables,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_Vegetables,
		},
		{
			Name: EventCreateFoodConfigurationFruits,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_Fruits,
		},
		{
			Name: EventCreateFoodConfigurationNutsAndSeeds,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_NutsAndSeeds,
		},
		{
			Name: EventCreateFoodConfigurationDairyProducts,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_DairyProducts,
		},
		{
			Name: EventCreateFoodConfigurationFish,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_Fish,
		},
		{
			Name: EventCreateFoodConfigurationMeat,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_Meat,
		},
		{
			Name: EventCreateFoodConfigurationEggs,
			Src:  []string{StateCreateFoodConfiguration},
			Dst:  StateCreateFoodConfiguration_Eggs,
		},
	}
}
