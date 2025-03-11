package entity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type DietConfiguration struct {
	ID                     string            `json:"id" bson:"_id"`
	UserID                 int64             `json:"user_id" bson:"user_id"`
	LifestyleAndWorkouts   string            `json:"lifestyle_and_workouts" bson:"lifestyle_and_workouts"`
	TimeRestrictions       string            `json:"time_restrictions" bson:"time_restrictions"`
	PFC                    PFC               `json:"pfc" bson:"pfc"`
	Calories               int               `json:"calories" bson:"calories"`
	NutritionPrinciples    string            `json:"nutrition_principles" bson:"nutrition_principles"`
	IndividualRestrictions string            `json:"individual_restrictions" bson:"individual_restrictions"`
	FoodConfiguration      FoodConfiguration `json:"food_configuration" bson:"food_configuration"`
}

func (d *DietConfiguration) CollectionName() string {
	return "diet_configurations"
}

type PFC struct {
	Proteins int `json:"proteins" bson:"proteins"`
	Fats     int `json:"fats" bson:"fats"`
	Carbs    int `json:"carbs" bson:"carbs"`
}

func (p *PFC) ParsePFC(pfc string) error {
	pfcPattern := `^\d{1,2}/\d{1,2}/\d{1,2}$`
	matched, err := regexp.MatchString(pfcPattern, pfc)
	if err != nil {
		return fmt.Errorf("error matching PFC pattern: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid PFC format")
	}

	parts := strings.Split(pfc, "/")
	if len(parts) != 3 {
		return fmt.Errorf("invalid PFC format")
	}

	proteins, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("error converting proteins to int: %w", err)
	}

	fats, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("error converting fats to int: %w", err)
	}

	carbs, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("error converting carbs to int: %w", err)
	}

	*p = PFC{
		Proteins: proteins,
		Fats:     fats,
		Carbs:    carbs,
	}

	return nil
}

type FoodConfiguration struct {
	BreakfastCereals string `json:"breakfast_cereals" bson:"breakfast_cereals"`
	SideDishCereals  string `json:"side_dish_cereals" bson:"side_dish_cereals"`
	Vegetables       string `json:"vegetables" bson:"vegetables"`
	Fruits           string `json:"fruits" bson:"fruits"`
	NutsAndSeeds     string `json:"nuts_and_seeds" bson:"nuts_and_seeds"`
	DairyProducts    string `json:"dairy_products" bson:"dairy_products"`
	Fish             string `json:"fish" bson:"fish"`
	Meat             string `json:"meat" bson:"meat"`
	Bread            string `json:"bread" bson:"bread"`
	Eggs             bool   `json:"eggs" bson:"eggs"`
}
