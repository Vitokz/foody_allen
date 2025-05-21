package flow

import "github.com/looplab/fsm"

const (
	EventUserConfigurationHeight    = "fill_user_configuration_height"
	EventUserConfigurationWeight    = "fill_user_configuration_weight"
	EventUserConfigurationGender    = "fill_user_configuration_gender"
	EventUserConfigurationAge       = "fill_user_configuration_age"
	EventUserConfigurationGoal      = "fill_user_configuration_goal"
	EventUserConfigurationActivity  = "fill_user_configuration_activity"
	EventUserConfigurationDietType  = "fill_user_configuration_diet_type"
	EventUserConfigurationAllergies = "fill_user_configuration_allergies"
	EventUserConfigurationMealTypes = "fill_user_configuration_meal_types"
)

func CreateUserConfigurationFlow() fsm.Events {
	return fsm.Events{
		{
			Name: EventUserConfigurationHeight,
			Src:  []string{StateMenu, StateStart},
			Dst:  StateUserConfiguration_Height,
		},
		{
			Name: EventUserConfigurationWeight,
			Src:  []string{StateMenu, StateUserConfiguration_Height},
			Dst:  StateUserConfiguration_Weight,
		},
		{
			Name: EventUserConfigurationGender,
			Src:  []string{StateMenu, StateUserConfiguration_Weight},
			Dst:  StateUserConfiguration_Gender,
		},
		{
			Name: EventUserConfigurationAge,
			Src:  []string{StateMenu, StateUserConfiguration_Gender},
			Dst:  StateUserConfiguration_Age,
		},
		{
			Name: EventUserConfigurationGoal,
			Src:  []string{StateMenu, StateUserConfiguration_Age},
			Dst:  StateUserConfiguration_Goal,
		},
		{
			Name: EventUserConfigurationActivity,
			Src:  []string{StateMenu, StateUserConfiguration_Goal},
			Dst:  StateUserConfiguration_Activity,
		},
		{
			Name: EventUserConfigurationDietType,
			Src:  []string{StateMenu, StateUserConfiguration_Activity},
			Dst:  StateUserConfiguration_DietType,
		},
		{
			Name: EventUserConfigurationAllergies,
			Src:  []string{StateMenu, StateUserConfiguration_DietType},
			Dst:  StateUserConfiguration_Allergies,
		},
		{
			Name: EventUserConfigurationMealTypes,
			Src:  []string{StateMenu, StateUserConfiguration_Allergies},
			Dst:  StateUserConfiguration_MealTypes,
		},
	}
}

func IsCreateUserConfigurationFillState(state string) bool {
	return state == StateUserConfiguration_Height ||
		state == StateUserConfiguration_Weight ||
		state == StateUserConfiguration_Gender ||
		state == StateUserConfiguration_Age ||
		state == StateUserConfiguration_Goal ||
		state == StateUserConfiguration_Activity ||
		state == StateUserConfiguration_DietType ||
		state == StateUserConfiguration_Allergies ||
		state == StateUserConfiguration_MealTypes
}
