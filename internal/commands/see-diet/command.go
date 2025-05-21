package seediet

import (
	"context"
	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
	"diet_bot/internal/repository"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DietCommand struct {
	logger     *zap.SugaredLogger
	repository *repository.Client
}

func NewDietCommand(repository *repository.Client, logger *zap.SugaredLogger) *DietCommand {
	return &DietCommand{
		logger:     logger,
		repository: repository,
	}
}

func (c *DietCommand) SeeDietHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	diet, err := c.repository.GetLatestDiet(meta.UserID)
	if err != nil {
		msg := tgbotapi.NewMessage(meta.ChatID, "Рацион еще не сформирован")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🥗 Создать рацион", flow.CommandGenerateDiet),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад в меню", flow.CommandMenu),
			),
		)

		return msg
	}

	msg := tgbotapi.NewMessage(meta.ChatID, "🥗 Выберите день рациона:")

	days := make([]tgbotapi.InlineKeyboardButton, len(diet.DailyDiet))
	for i := range diet.DailyDiet {
		dayNumber := strconv.Itoa(i + 1)
		days[i] = tgbotapi.NewInlineKeyboardButtonData("🥗 "+dayNumber, dietDayCommand(dayNumber))
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, days)
	backButton := tgbotapi.NewInlineKeyboardButtonData("🔙 Назад в меню", flow.CommandMenu)
	rows = append(rows, []tgbotapi.InlineKeyboardButton{backButton})

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg
}

func (c *DietCommand) SeeDietDayHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	dayNumber, err := getDietDayFromCommand(update.CallbackQuery.Data)
	if err != nil {
		c.logger.Error("error getting diet day", zap.Error(err))
		return nil
	}

	diet, err := c.repository.GetLatestDiet(meta.UserID)
	if err != nil {
		c.logger.Error("error getting diet", zap.Error(err))
		return nil
	}

	day := diet.DailyDiet[dayNumber-1]
	msg := tgbotapi.NewMessage(meta.ChatID, day.ToMessage())
	msg.ParseMode = tgbotapi.ModeMarkdown

	backButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад к списку", flow.CommandSeeDiet),
		),
	)

	msg.ReplyMarkup = backButton

	return msg
}

func (c *DietCommand) SeeDietProductsHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	diet, err := c.repository.GetLatestDiet(meta.UserID)
	if err != nil && err != mongo.ErrNoDocuments {
		c.logger.Error("error getting diet", zap.Error(err))
		return nil
	}

	if err == mongo.ErrNoDocuments {
		msg := tgbotapi.NewMessage(meta.ChatID, "Рацион еще не сформирован")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🥗 Создать рацион", flow.CommandGenerateDiet),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад в меню", flow.CommandMenu),
			),
		)

		return msg
	}

	// Создаем карту для группировки продуктов по категориям
	productsByType := make(map[entity.ProductToBuyType][]entity.ProductToBuy)

	// Заполняем карту продуктами
	for _, product := range diet.ProductsToBuy {
		productType := entity.ProductToBuyType(product.Type)
		productsByType[productType] = append(productsByType[productType], product)
	}

	// Формируем сообщение, следуя порядку категорий из ProductToBuySort
	var messageBuilder strings.Builder
	messageBuilder.WriteString("🛒 *Список продуктов для покупки:*\n\n")

	totalProducts := 0 // Счетчик всех продуктов

	// Проходим по категориям в нужном порядке
	for _, category := range entity.ProductToBuySort {
		products := productsByType[category]
		if len(products) == 0 {
			continue // Пропускаем пустые категории
		}

		// Добавляем заголовок категории с соответствующим эмодзи
		emoji := getCategoryEmoji(category)
		messageBuilder.WriteString(fmt.Sprintf("*%s %s:*\n", emoji, category))

		// Добавляем продукты категории
		for _, product := range products {
			weight := product.Weight
			unit := "г"

			// Конвертируем большие веса в кг для удобства чтения
			if weight >= 1000 {
				weight = weight / 1000
				unit = "кг"
			}

			messageBuilder.WriteString(fmt.Sprintf("• %s — %d%s\n", product.Name, weight, unit))
			totalProducts++
		}

		messageBuilder.WriteString("\n")
	}

	// Добавляем итоговую информацию
	messageBuilder.WriteString(fmt.Sprintf("*Всего продуктов:* %d\n", totalProducts))

	// Создаем сообщение с Markdown форматированием
	msg := tgbotapi.NewMessage(meta.ChatID, messageBuilder.String())
	msg.ParseMode = tgbotapi.ModeMarkdown

	// Добавляем кнопку "Назад в меню"
	backButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад в меню", flow.CommandMenu),
		),
	)
	msg.ReplyMarkup = backButton

	return msg
}

// Вспомогательная функция для получения эмодзи категории
func getCategoryEmoji(category entity.ProductToBuyType) string {
	switch category {
	case entity.ProductToBuyTypeVegetables:
		return "🥦"
	case entity.ProductToBuyTypeMilk:
		return "🥛"
	case entity.ProductToBuyTypeMeatAndFish:
		return "🥩"
	case entity.ProductToBuyTypeCereals:
		return "🌾"
	case entity.ProductToBuyTypeAdditional:
		return "🧂"
	default:
		return "📋"
	}
}

func dietDayCommand(dayNumber string) string {
	return flow.CommandSeeDietDay + dayNumber
}

func getDietDayFromCommand(command string) (int, error) {
	parts := strings.Split(command, "_")

	dayNumber, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, err
	}

	return dayNumber, nil
}

func CommandHasDietDay(command string) bool {
	return strings.HasPrefix(command, flow.CommandSeeDietDay)
}
