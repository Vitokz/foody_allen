package entity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type DietConfiguration struct {
	ID                     uuid.UUID         `json:"id" bson:"_id"`
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
	Proteins float64 `json:"proteins" bson:"proteins" jsonschema_description:"Количество белков"`
	Fats     float64 `json:"fats" bson:"fats" jsonschema_description:"Количество жиров"`
	Carbs    float64 `json:"carbs" bson:"carbs" jsonschema_description:"Количество углеводов"`
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

	proteins, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("error converting proteins to float64: %w", err)
	}

	fats, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("error converting fats to float64: %w", err)
	}

	carbs, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return fmt.Errorf("error converting carbs to float64: %w", err)
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
