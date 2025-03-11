package flow

import "github.com/looplab/fsm"

const (
	EventCreateDiet                       = "create_diet"
	EventCreateDietTimeRestrictions       = "create_diet_time_restrictions"
	EventCreateDietPFC                    = "create_diet_pfc"
	EventCreateDietCalories               = "create_diet_calories"
	EventCreateDietNutritionPrinciples    = "create_diet_nutrition_principles"
	EventCreateDietIndividualRestrictions = "create_diet_individual_restrictions"
)

func CreateDietFlow() fsm.Events {
	return fsm.Events{
		{
			Name: EventCreateDiet,
			Src:  []string{StateMenu, StateStart},
			Dst:  StateCreateDiet_LifeStyle,
		},
		{
			Name: EventCreateDietTimeRestrictions,
			Src:  []string{StateCreateDiet_LifeStyle},
			Dst:  StateCreateDiet_TimeRestrictions,
		},
		{
			Name: EventCreateDietPFC,
			Src:  []string{StateCreateDiet_TimeRestrictions},
			Dst:  StateCreateDiet_PFC,
		},
		{
			Name: EventCreateDietCalories,
			Src:  []string{StateCreateDiet_PFC},
			Dst:  StateCreateDiet_Calories,
		},
		{
			Name: EventCreateDietNutritionPrinciples,
			Src:  []string{StateCreateDiet_Calories},
			Dst:  StateCreateDiet_NutritionPrinciples,
		},
		{
			Name: EventCreateDietIndividualRestrictions,
			Src:  []string{StateCreateDiet_NutritionPrinciples},
			Dst:  StateCreateDiet_IndividualRestrictions,
		},
	}
}

func IsCreateDietFillState(state string) bool {
	return state == StateCreateDiet_LifeStyle ||
		state == StateCreateDiet_TimeRestrictions ||
		state == StateCreateDiet_PFC ||
		state == StateCreateDiet_Calories ||
		state == StateCreateDiet_NutritionPrinciples ||
		state == StateCreateDiet_IndividualRestrictions
}


