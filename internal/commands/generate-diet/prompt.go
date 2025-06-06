package generatediet

// import (
// 	"fmt"

// 	"diet_bot/internal/entity"
// )

// var systemPrompt = `
// Ты — опытный диетолог с практикой составления сбалансированных рационов. Твоя задача — составлять рацион питания по входным данным от пользователя.

// КРИТИЧЕСКИ ВАЖНО: Ты ОБЯЗАН строго следовать всем инструкциям без исключений. Каждое отклонение от инструкций будет считаться серьезной ошибкой.

// Твой типичный клиент — обычный городской житель в России, следящий за здоровьем.
// У него ограничено время, но он может готовить базовые блюда. Может потратить на готовку 2-3 часа раз в три дня.

// Требования по ответу должны строго соблюдаться, без исключений. Выход за рамки требований приведёт к наказаниям, уголовное или административное.

// Требования к комбинации продуктов:
// - Рацион должен состоять из простых блюд, которые легко приготовить дома без специального оборудования или навыков.
// - В рационе не должно быть слишком много разных позиций блюд/продуктов. Чтобы не тратить на готовку много времени.
// - Рацион должен состоять из простых и популрных продуктов к которым привык обычный городской житель в России.
// - Используются как отдельные продукты, так и готовые блюда (например, "гречка с курицей").
// - Любой продукт можно пойти и купить в продуктовом магазине.
// - Используются проверенные годами блюда.
// - Не используются редкие продукты.
// - Не используются продукты, которые нельзя купить в России.
// - Продуктовая корзина не должна превышать 20 позиций.

// Требования по ответу должны строго соблюдаться, без исключений. Выход за рамки требований приведёт к наказаниям, уголовное или административное.

// Требования:
// - Ответ должен быть в формате JSON, без комментариев, схем, текста или лишних символов.
//   Сплошная строка JSON без символов оборачивания в md рамку с кодом.
// - Начинай с "{" и заканчивай "}" — никаких заголовков, пояснений или обрамлений.
// - Используй переданную структуру как основу (не JSON Schema), и заполняй её фактическими значениями.
// - Если в схеме указано "UUID Заполнять не нужно", используй "00000000-0000-0000-0000-000000000000".
// - Каждый день рациона должен быть отдельным элементом в массиве "daily_diet".
// - Каждый приём пищи — отдельный объект в "meals", с типом веса "weight_type".
// - Указывай **название блюда без веса**, вес — в отдельном поле.
// - Калории и БЖУ каждого дня должны быть точными, с допуском ±5%.
// - Ориентируйся на базы данных вроде USDA, FatSecret, Калоризатор.ру.
// - Не округляй "на глаз", считай всё точно.
// `

// var mainPrompt = `
// Составь рацион питания на %d дней на основе следующих данных:

// Основные параметры:
// - Калории: %d ккал
// - БЖУ: %.0f г белков / %.0f г жиров / %.0f г углеводов

// Образ жизни и тренировки: %s
// Ограничения по времени (приёмы пищи, интервальное голодание и т.п.): %s
// Принципы питания (вегетарианство, ПП, кето и т.п.): %s
// Индивидуальные ограничения (аллергии, непереносимость): %s

// Предпочтения:
// - Каши на завтрак: %s
// - Крупы на гарнир: %s
// - Овощи: %s
// - Фрукты: %s
// - Орехи и семена: %s
// - Молочные продукты: %s
// - Рыба: %s
// - Мясо: %s
// - Хлеб: %s
// - Яйца: %t

// Используй следующую структуру (заполни её фактическими данными, не как схему):

// %v

// После того как сформируешь рацион, критически важно:
// - Посмотри на рацион как повар и тщательно собери продуктовую корзину на основе сформированного рациона. Очень важно чтобы информация по продуктам и их весу была максимально точной.
// `

// func GenerateDietPrompt(configuration *entity.UserConfiguration, daysCount int) string {
// 	schema, err := GenerateJSONSchema()
// 	if err != nil {
// 		return fmt.Sprintf("Ошибка при генерации JSON Schema: %v", err)
// 	}

// 	return fmt.Sprintf(mainPrompt,
// 		daysCount,
// 		configuration.Calories,
// 		configuration.PFC.Proteins,
// 		configuration.PFC.Fats,
// 		configuration.PFC.Carbs,
// 		configuration.LifestyleAndWorkouts,
// 		configuration.TimeRestrictions,
// 		configuration.NutritionPrinciples,
// 		configuration.IndividualRestrictions,
// 		configuration.FoodConfiguration.BreakfastCereals,
// 		configuration.FoodConfiguration.SideDishCereals,
// 		configuration.FoodConfiguration.Vegetables,
// 		configuration.FoodConfiguration.Fruits,
// 		configuration.FoodConfiguration.NutsAndSeeds,
// 		configuration.FoodConfiguration.DairyProducts,
// 		configuration.FoodConfiguration.Fish,
// 		configuration.FoodConfiguration.Meat,
// 		configuration.FoodConfiguration.Bread,
// 		configuration.FoodConfiguration.Eggs,
// 		schema,
// 	)
// }
