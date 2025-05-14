package generatediet

import (
	"fmt"

	"diet_bot/internal/entity"
)

var mainPrompt = `
	Я хочу составить рацион питания на %d дней. Войди в роль диетолога и составь рацион питания на основе следующих данных:

	Образ жизни и тренировки: %s
	Ограничения по времени: %s
	БЖУ: %.0f/%.0f/%.0f
	Калории: %d
	Принципы питания: %s
	Индивидуальные ограничения: %s

	Предпочтения в еде:
	- Каши на завтрак: %s
	- Крупы на гарнир: %s
	- Овощи: %s
	- Фрукты: %s
	- Орехи и семена: %s
	- Молочные продукты: %s
	- Рыба: %s
	- Мясо: %s
	- Хлеб: %s
	- Яйца: %t

	Вот требования к ответу:
	- В ответе отдай только JSON для рациона питания без лишних символов, чтобы я мог его сразу спарсить
	- В ответе отдавай не только отдельные продукты, но и целые блюда. Например "гречка с курицей 100г"
	- В названии не пиши вес продукта. Для этого есть отдельное поле.
	- Обязательно проследи чтобы количество калорий во всех приемах пищи совпадало с общим количеством калорий в рационе.
	- Если в схеме ответа в description к полю указано "UUID Заполнять не нужно", то подставь 00000000-0000-0000-0000-000000000000.
	- Не нужно подставлять какие-либо символы. Только структура JSON не более. Если подставишь что-то лишнее я подам на тебя в суд
	- Не отдавай мне JSON Schema, мне нужна готовая структура JSON, чтобы я мог ее подставить в JS код сразу.
	Вот схема ответа: 
	%v
`

func GenerateDietPrompt(configuration *entity.DietConfiguration, daysCount int) string {
	schema, err := GenerateJSONSchema()
	if err != nil {
		return fmt.Sprintf("Ошибка при генерации JSON Schema: %v", err)
	}

	return fmt.Sprintf(mainPrompt,
		daysCount,
		configuration.LifestyleAndWorkouts,
		configuration.TimeRestrictions,
		configuration.PFC.Proteins,
		configuration.PFC.Fats,
		configuration.PFC.Carbs,
		configuration.Calories,
		configuration.NutritionPrinciples,
		configuration.IndividualRestrictions,
		configuration.FoodConfiguration.BreakfastCereals,
		configuration.FoodConfiguration.SideDishCereals,
		configuration.FoodConfiguration.Vegetables,
		configuration.FoodConfiguration.Fruits,
		configuration.FoodConfiguration.NutsAndSeeds,
		configuration.FoodConfiguration.DairyProducts,
		configuration.FoodConfiguration.Fish,
		configuration.FoodConfiguration.Meat,
		configuration.FoodConfiguration.Bread,
		configuration.FoodConfiguration.Eggs,
		schema,
	)
}
