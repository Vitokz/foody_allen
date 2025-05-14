package generatediet

import (
	"fmt"

	"diet_bot/internal/entity"
)

var mainPrompt = `
Ты — опытный диетолог с практикой составления сбалансированных рационов. Составь рацион питания на %d дней на основе следующих данных:

Основные параметры:
- Калории: %d ккал
- БЖУ: %.0f г белков / %.0f г жиров / %.0f г углеводов

Образ жизни и тренировки: %s  
Ограничения по времени (приёмы пищи, интервальное голодание и т.п.): %s  
Принципы питания (вегетарианство, ПП, кето и т.п.): %s  
Индивидуальные ограничения (аллергии, непереносимость): %s

Предпочтения:
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

Требования к результату:
- Верни только JSON — **без текста до или после**, без схем, только чистый JSON, сразу готовый к парсингу в коде.
- Структура должна соответствовать следующей схеме:  
  %v
- Если в схеме указано "UUID Заполнять не нужно", то используй "00000000-0000-0000-0000-000000000000"
- Используй как отдельные продукты, так и готовые блюда (например, "гречка с курицей")
- Вес блюда указывай в отдельном поле, не в названии
- Обязательно соблюдай баланс: суммарная калорийность за день должна строго соответствовать указанной (±5%), так же как и БЖУ (±5%)
`

func GenerateDietPrompt(configuration *entity.DietConfiguration, daysCount int) string {
	schema, err := GenerateJSONSchema()
	if err != nil {
		return fmt.Sprintf("Ошибка при генерации JSON Schema: %v", err)
	}

	return fmt.Sprintf(mainPrompt,
		daysCount,
		configuration.Calories,
		configuration.PFC.Proteins,
		configuration.PFC.Fats,
		configuration.PFC.Carbs,
		configuration.LifestyleAndWorkouts,
		configuration.TimeRestrictions,
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
