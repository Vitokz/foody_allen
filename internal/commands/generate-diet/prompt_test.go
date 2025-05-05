package generatediet

import (
	"diet_bot/internal/entity"
	"fmt"
	"testing"
)

func TestGenerateSchemaPrompt(t *testing.T) {
	schema, err := GenerateJSONSchema()
	if err != nil {
		t.Fatalf("Failed to generate schema: %v", err)
	}

	fmt.Println(schema)
}

func TestGenerateDietPrompt(t *testing.T) {
	configuration := &entity.DietConfiguration{
		LifestyleAndWorkouts: "Сидячая работа, тренируюсь 4 раза в неделю по 2-3 часа кроссфит",
		TimeRestrictions:     "Нет ограничений по времени",
		PFC: entity.PFC{
			Proteins: 2.5,
			Fats:     1.5,
			Carbs:    3.0,
		},
		Calories:               2000,
		NutritionPrinciples:    "Нет принципов питания",
		IndividualRestrictions: "Нет индивидуальных ограничений",
		FoodConfiguration: entity.FoodConfiguration{
			BreakfastCereals: "Рис",
			SideDishCereals:  "Рис",
			Vegetables:       "Овощи",
			Fruits:           "Фрукты",
			NutsAndSeeds:     "Орехи и семена",
			DairyProducts:    "Молочные продукты",
			Fish:             "Рыба",
			Meat:             "Мясо",
			Bread:            "Хлеб",
			Eggs:             false,
		},
	}

	prompt := GenerateDietPrompt(configuration, 3)
	fmt.Println(prompt)
}
