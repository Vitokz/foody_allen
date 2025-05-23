package generatediet

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
	"diet_bot/internal/repository"
)

type AIClient interface {
	GenerateDiet(systemPrompt string, prompt string) (string, error)
}

type Command struct {
	repository *repository.Client
	aiClient   AIClient
	logger     *zap.SugaredLogger
}

func NewCommand(repository *repository.Client, aiClient AIClient, logger *zap.SugaredLogger) *Command {
	return &Command{
		repository: repository,
		aiClient:   aiClient,
		logger:     logger,
	}
}

func (c *Command) GenerateDietHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	msg := tgbotapi.NewMessage(meta.ChatID, "Выбери количество дней для рациона:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1 день", generateDietDayCommand("1")),
			tgbotapi.NewInlineKeyboardButtonData("3 дня", generateDietDayCommand("3")),
			tgbotapi.NewInlineKeyboardButtonData("5 дней", generateDietDayCommand("5")),
			tgbotapi.NewInlineKeyboardButtonData("7 дней", generateDietDayCommand("7")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад в меню", flow.CommandMenu),
		),
	)

	return msg
}

func generateDietDayCommand(dayNumber string) string {
	return flow.CommandGenerateDietDays + dayNumber
}

func generateDietDayFromCommand(command string) (int, error) {
	parts := strings.Split(command, "_")

	dayNumber, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, err
	}

	return dayNumber, nil
}

func CommandHasGenerateDietDays(command string) bool {
	return strings.HasPrefix(command, flow.CommandGenerateDietDays)
}

func (c *Command) GenerateDietDaysHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	_, err := generateDietDayFromCommand(update.CallbackQuery.Data)
	if err != nil {
		c.logger.Error("error getting diet days count", zap.Error(err))
		return nil
	}

	response, err := c.aiClient.GenerateDiet("", "")
	if err != nil {
		c.logger.Error("error generating diet", zap.Error(err))
		return nil
	}

	c.logger.Info("generated diet", zap.String("diet", response))

	diet := new(entity.GeneratedDiet)
	err = json.Unmarshal([]byte(response), diet)
	if err != nil {
		c.logger.Error("error unmarshalling diet", zap.Error(err))
		return nil
	}

	diet.SetIDs()
	diet.UserID = meta.UserID
	diet.CreatedAt = time.Now()

	err = c.repository.CreateDiet(diet)
	if err != nil {
		c.logger.Error("error creating diet", zap.Error(err))
		return nil
	}

	msg := tgbotapi.NewMessage(meta.ChatID, "Рацион успешно создан!")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥗 Посмотреть рацион", flow.CommandSeeDiet),
		),
	)

	return msg
}
