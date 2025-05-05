package entity

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type GeneratedDiet struct {
	ID            uuid.UUID      `json:"id" bson:"_id" jsonschema:"type=string" jsonschema_description:"UUID Заполнять не нужно"`
	UserID        int64          `json:"user_id" bson:"user_id" jsonschema_description:"Заполнять не нужно"`
	ConfigID      uuid.UUID      `json:"config_id" bson:"config_id" jsonschema:"type=string" jsonschema_description:"UUID Заполнять не нужно"`
	DailyDiet     []DailyDiet    `json:"daily_diet" bson:"daily_diet" jsonschema_description:"Ежедневный рацион"`
	ProductsToBuy []ProductToBuy `json:"products_to_buy" bson:"products_to_buy" jsonschema_description:"Продуктовая корзина"`
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
	Name     string `json:"name" bson:"name" jsonschema_description:"Название приема пищи: например"`
	Time     string `json:"time" bson:"time" jsonschema_description:"Время приема пищи"`
	Calories int    `json:"calories" bson:"calories" jsonschema_description:"Количество калорий в приеме пищи"`
	PFC      PFC    `json:"pfc" bson:"pfc" jsonschema_description:"Количество белков, жиров и углеводов в приеме пищи"`
	Dishes   []Dish `json:"dishes" bson:"dishes" jsonschema_description:"Что входит в прием пищи"`
}

type DishWeightType string

const (
	DishWeightTypeDry    DishWeightType = "сухой"
	DishWeightTypeCooked DishWeightType = "приготовленный"
)

type Dish struct {
	Name        string         `json:"name" bson:"name" jsonschema_description:"Название продукта"`
	Calories    int            `json:"calories" bson:"calories" jsonschema_description:"Количество калорий в продукте"`
	PFC         PFC            `json:"pfc" bson:"pfc" jsonschema_description:"Количество белков, жиров и углеводов в продукте"`
	Weight      int            `json:"weight" bson:"weight" jsonschema_description:"Количество грамм в продукте"`
	WeightUnits string         `json:"weight_units" bson:"weight_units" jsonschema_description:"Единицы измерения веса в продукте"`
	WeightType  DishWeightType `json:"weight_type" bson:"weight_type" jsonschema_description:"Тип веса в продукте" jsonschema:"enum=сухой,enum=приготовленный"`
}

type ProductToBuyType string

const (
	ProductToBuyTypeVegetables  ProductToBuyType = "Овощи и фрукты"
	ProductToBuyTypeMilk        ProductToBuyType = "Молочка"
	ProductToBuyTypeMeatAndFish ProductToBuyType = "Мясо и рыба"
	ProductToBuyTypeCereals     ProductToBuyType = "Крупы"
	ProductToBuyTypeAdditional  ProductToBuyType = "Дополнительно"
)

var ProductToBuySort = []ProductToBuyType{
	ProductToBuyTypeVegetables,
	ProductToBuyTypeMilk,
	ProductToBuyTypeMeatAndFish,
	ProductToBuyTypeCereals,
	ProductToBuyTypeAdditional,
}

type ProductToBuy struct {
	Name   string `json:"name" bson:"name" jsonschema_description:"Название продукта"`
	Weight int    `json:"weight" bson:"weight" jsonschema_description:"Количество грамм в продукте"`
	Type   string `json:"type" bson:"type" jsonschema_description:"Тип продукта" jsonschema:"enum=Овощи и фрукты,enum=Молочка,enum=Мясо и рыба,enum=Крупы,enum=Дополнительно"`
}

func (d *DailyDiet) ToMessage() string {
	message := "🍽️ *Ваш рацион на день* 🍽️\n\n"
	message += fmt.Sprintf("🔥 *Общая калорийность:* %d ккал\n", d.TotalCalories)
	message += fmt.Sprintf("🥩 *Белки:* %.1f г\n", d.TotalPFC.Proteins)
	message += fmt.Sprintf("🧈 *Жиры:* %.1f г\n", d.TotalPFC.Fats)
	message += fmt.Sprintf("🍚 *Углеводы:* %.1f г\n\n", d.TotalPFC.Carbs)

	for i, meal := range d.Meals {
		timeEmoji := timeToNumberEmoji(meal.Time)
		message += fmt.Sprintf("*%s %s - %d ккал\n", meal.Name, timeEmoji, meal.Calories)
		message += fmt.Sprintf("🥩 %.1fг, 🧈 %.1fг, 🍚 %.1fг\n\n",
			meal.PFC.Proteins, meal.PFC.Fats, meal.PFC.Carbs)

		message += "*Продукты:*\n"
		for _, dish := range meal.Dishes {
			weightTypeInfo := ""
			if dish.WeightType == DishWeightTypeDry {
				weightTypeInfo = " (сухой)"
			}

			message += fmt.Sprintf("🍲 %s - %d %s%s (%d ккал)\n",
				dish.Name, dish.Weight, dish.WeightUnits, weightTypeInfo, dish.Calories)
		}

		if i < len(d.Meals)-1 {
			message += "\n" + strings.Repeat("—", 20) + "\n\n"
		}
	}

	return message
}

func timeToNumberEmoji(timeStr string) string {
	emojiMap := map[rune]string{
		'0': "0️⃣",
		'1': "1️⃣",
		'2': "2️⃣",
		'3': "3️⃣",
		'4': "4️⃣",
		'5': "5️⃣",
		'6': "6️⃣",
		'7': "7️⃣",
		'8': "8️⃣",
		'9': "9️⃣",
		':': ":",
	}

	result := ""
	for _, char := range timeStr {
		if emoji, exists := emojiMap[char]; exists {
			result += emoji
		} else {
			result += string(char)
		}
	}

	return result
}
