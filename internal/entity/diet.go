package entity

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type GeneratedDiet struct {
	ID        uuid.UUID   `json:"id" bson:"_id" jsonschema:"type=string" jsonschema_description:"UUID Заполнять не нужно"`
	UserID    int64       `json:"user_id" bson:"user_id" jsonschema_description:"Заполнять не нужно"`
	ConfigID  uuid.UUID   `json:"config_id" bson:"config_id" jsonschema:"type=string" jsonschema_description:"UUID Заполнять не нужно"`
	DailyDiet []DailyDiet `json:"daily_diet" bson:"daily_diet" jsonschema_description:"Ежедневный рацион"`
}

func (d *GeneratedDiet) CollectionName() string {
	return "generated_diets"
}

func (d *GeneratedDiet) SetIDs() {
	d.ID = uuid.New()

	for dailyIdx, dailyDiet := range d.DailyDiet {
		dailyDiet.ID = uuid.New()

		d.DailyDiet[dailyIdx] = dailyDiet
	}
}

type DailyDiet struct {
	ID            uuid.UUID `json:"id" bson:"_id,omitempty" jsonschema:"type=string" jsonschema_description:"UUID Заполнять не нужно"`
	TotalCalories int       `json:"total_calories" bson:"total_calories" jsonschema_description:"Общее количество калорий за день"`
	TotalPFC      PFC       `json:"total_pfc" bson:"total_pfc" jsonschema_description:"Общее количество белков, жиров и углеводов за день"`
	Meals         []Meal    `json:"meals" bson:"meals" jsonschema_description:"Приемы пищи"`
}

type Meal struct {
	Name     string    `json:"name" bson:"name" jsonschema_description:"Название приема пищи: например"`
	Time     string    `json:"time" bson:"time" jsonschema_description:"Время приема пищи"`
	Calories int       `json:"calories" bson:"calories" jsonschema_description:"Количество калорий в приеме пищи"`
	PFC      PFC       `json:"pfc" bson:"pfc" jsonschema_description:"Количество белков, жиров и углеводов в приеме пищи"`
	Products []Product `json:"products" bson:"products" jsonschema_description:"Продукты в приеме пищи"`
}

type Product struct {
	Name     string `json:"name" bson:"name" jsonschema_description:"Название продукта"`
	Calories int    `json:"calories" bson:"calories" jsonschema_description:"Количество калорий в продукте"`
	PFC      PFC    `json:"pfc" bson:"pfc" jsonschema_description:"Количество белков, жиров и углеводов в продукте"`
	Weight   int    `json:"weight" bson:"weight" jsonschema_description:"Количество грамм в продукте"`
}

func (d *DailyDiet) ToMessage() string {
	message := "🍽️ *Ваш рацион на день* 🍽️\n\n"
	message += fmt.Sprintf("🔥 *Общая калорийность:* %d ккал\n", d.TotalCalories)
	message += fmt.Sprintf("🥩 *Белки:* %.1f г\n", d.TotalPFC.Proteins)
	message += fmt.Sprintf("🧈 *Жиры:* %.1f г\n", d.TotalPFC.Fats)
	message += fmt.Sprintf("🍚 *Углеводы:* %.1f г\n\n", d.TotalPFC.Carbs)

	for i, meal := range d.Meals {
		message += fmt.Sprintf("*%s (%s)* - %d ккал\n", meal.Name, meal.Time, meal.Calories)
		message += fmt.Sprintf("Б: %.1fг, Ж: %.1fг, У: %.1fг\n\n",
			meal.PFC.Proteins, meal.PFC.Fats, meal.PFC.Carbs)

		message += "*Продукты:*\n"
		for j, product := range meal.Products {
			message += fmt.Sprintf("%d. %s - %dг (%d ккал)\n",
				j+1, product.Name, product.Weight, product.Calories)
		}

		if i < len(d.Meals)-1 {
			message += "\n" + strings.Repeat("—", 20) + "\n\n"
		}
	}

	return message
}
