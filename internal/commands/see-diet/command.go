package seediet

import (
	"context"
	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
	"strconv"
	"strings"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Repository interface {
	GetDiet(chatID int64) (*entity.GeneratedDiet, error)
}

type DietCommand struct {
	logger     *zap.SugaredLogger
	repository Repository
}

func NewDietCommand(repository Repository, logger *zap.SugaredLogger) *DietCommand {
	return &DietCommand{
		logger:     logger,
		repository: repository,
	}
}

func (c *DietCommand) SeeDietHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	diet, err := c.repository.GetDiet(meta.ChatID)
	if err != nil {
		c.logger.Error("error getting diet", zap.Error(err))
		return nil
	}

	msg := tgbotapi.NewMessage(meta.ChatID, "🥗 Выберите день рациона:")

	days := make([]tgbotapi.InlineKeyboardButton, len(diet.DailyDiet))
	for i := range diet.DailyDiet {
		dayNumber := strconv.Itoa(i + 1)

		days[i] = tgbotapi.NewInlineKeyboardButtonData("🥗 "+dayNumber, dietDayCommand(dayNumber))
	}

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			days...,
		),
	)

	return msg
}

func (c *DietCommand) SeeDietDayHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	dayNumber, err := getDietDayFromCommand(update.CallbackQuery.Data)
	if err != nil {
		c.logger.Error("error getting diet day", zap.Error(err))
		return nil
	}

	diet, err := c.repository.GetDiet(meta.ChatID) // Здесь я использую chatID, потому что в этом месте я получаю callbackQuery, а не message и юзер в данном случае это бот
	if err != nil {
		c.logger.Error("error getting diet", zap.Error(err))
		return nil
	}

	day := diet.DailyDiet[dayNumber-1]

	msg := tgbotapi.NewMessage(meta.ChatID, day.ToMessage())
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
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
