package creatediet

import (
	"context"

	"diet_bot/internal/entity"
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Repository interface {
	GetChat(id int64) (*entity.Chat, error)
	UpsertChat(chat *entity.Chat) error

	GetDietConfiguration(userID int64) (*entity.DietConfiguration, error)
	UpsertDietConfiguration(dietConfiguration *entity.DietConfiguration) error
}

type Commands struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func NewCommands(repository Repository, logger *zap.SugaredLogger) *Commands {
	return &Commands{
		repository: repository,
		logger:     logger,
	}
}

func (c *Commands) IsFillDiet(ctx context.Context, update *tgbotapi.Update) bool {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return false
	}

	ok := chat.State == flow.StateCreateDiet_LifeStyle

	if !ok {
		c.logger.Infof("Chat state: %s is not fill diet", chat.State)
	} else {
		c.logger.Infof("Chat state: %s is fill diet", chat.State)
	}

	return ok
}

func (c *Commands) CreateDietHandler(ctx context.Context, update *tgbotapi.Update) tgbotapi.Chattable {
	text := `
❔ Давай начнем с того, какой у тебя образ жизни?

Опиши свой образ жизни своими словами в свободной форме.

Например: Работаю в офисе, сижу за компьютером. В день прохожу примерно 3 км. 
Занимаюсь кроссфитом 4 раза в неделю, по 2-3 часа каждая тренировка.
`

	meta := entity.NewMeta(update)

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка нажмите /start еще раз")
	}

	botFSM := flow.NewBotFSM(chat)

	if err := botFSM.Event(flow.EventCreateDiet); err != nil {
		c.logger.Errorw("Error creating diet", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при создании диеты")
	}

	if err := c.repository.UpsertChat(chat); err != nil {
		c.logger.Errorw("Error upserting chat", "error", err)
		return tgbotapi.NewMessage(meta.ChatID, "Произошла ошибка при сохранении данных")
	}

	return msg
}

func (c *Commands) TimeRestrictionsHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Какие у тебя есть ограничения по времени?

Например: Тренировка с 18:00 до 21:00, по-этому нет возможности есть в это время.
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateDietTimeRestrictions); err != nil {
		c.logger.Errorw("Error creating diet time restrictions", "error", err)

		return nil, err
	}

	return msg, nil
}
