package entity

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func DefaultDietConfiguration() *DietConfiguration {
	return &DietConfiguration{
		LifestyleAndWorkouts: `Работаю в офисе, сижу за компьютером. В день прохожу примерно 5-7 км. 
Занимаюсь фитнесом 2-3 раза в неделю по 1-1.5 часа.`,
		TimeRestrictions: "Рабочий день с 9:00 до 18:00. Обед с 13:00 до 14:00. Ужин не позже 20:00.",
		PFC: PFC{
			Proteins: 20,
			Fats:     30,
			Carbs:    40,
		},
		Calories: 2400,
		NutritionPrinciples: `1. Без перекусов.
2. 3 приема пищи.
3. Разнообразный рацион.
4. Упор на медленные углеводы и овощи.
5. Завтрак - самый важный прием пищи.`,
		IndividualRestrictions: `1. Не есть после 21:00
2. Ограничить сладкое до 1-2 раз в неделю
3. Должно быть 2-3 вида гарнира на неделю
4. Должно быть 2 вида мяса на неделю
5. Должно быть 2 вида рыбы на неделю
6. Минимум 3 разных вида овощей в день`,
		FoodConfiguration: FoodConfiguration{
			BreakfastCereals: "овсянка,мюсли,гречневая каша,кукурузные хлопья",
			SideDishCereals:  "рис,гречка,макароны,картофель,булгур,киноа",
			Vegetables:       "огурцы,помидоры,капуста,морковь,лук,перец,баклажаны,кабачки,брокколи,шпинат,руккола,авокадо",
			Fruits:           "яблоко,банан,апельсин,киви,виноград,клубника,черника",
			NutsAndSeeds:     "миндаль,кешью,грецкий орех,семечки подсолнуха,тыквенные семечки",
			DairyProducts:    "молоко,кефир,йогурт,творог,сыр,сметана",
			Fish:             "лосось,тунец,минтай,треска,сельдь",
			Meat:             "курица,индейка,говядина,свинина",
			Bread:            "цельнозерновой,ржаной,багет,лаваш",
			Eggs:             true,
		},
	}
}

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
