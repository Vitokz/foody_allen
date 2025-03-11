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

func (c *Commands) PFCHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
Введи количество белков, жиров и углеводов в процентах. 

Формат: 30/30/40

Например: 20/30/50
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateDietPFC); err != nil {
		c.logger.Errorw("Error creating diet pfc", "error", err)

		return nil, err
	}

	return msg, nil
}

func (c *Commands) CaloriesHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Какое количество калорий тебе нужно в день?

Формат: 2000

Например: 3250
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateDietCalories); err != nil {
		c.logger.Errorw("Error creating diet calories", "error", err)
	}

	return msg, nil
}

func (c *Commands) NutritionPrinciplesHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Какие принципы питания ты хочешь поддерживать?

Формат свободный, но чтобы было понятно, что это за принципы.

Например: 
1. Без перекусов.
2. Дробное питание на 5 приемов.
3. Разнообразие рациона.
4. Упор на медленные углеводы и овощи.
5. Завтрак - самый главный прием пищи.
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateDietNutritionPrinciples); err != nil {
		c.logger.Errorw("Error creating diet nutrition principles", "error", err)
	}

	return msg, nil
}

func (c *Commands) IndividualRestrictionsHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
❔ Какие у тебя есть ограничения?

Формат свободный, но чтобы было понятно, что это за ограничения.

Например: 
1. Не есть после 18:00.
2. Не есть сладкое.
3. Не использовать бобовые
4. Не использовать зеленый лук и болгарский перец.
5. Должно быть два вида гарнира, чтобы не приходилось готовить слишком много
6. Должен быть 1 тип рыбы. Чтобы не готовить слишком много
7. Должен быть 1 тип мяса. Чтобы не готовить слишком много
8. Два вида тушеных овощей вместо множества отдельных блюд. Чтобы не готовить слишком много.
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateDietIndividualRestrictions); err != nil {
		c.logger.Errorw("Error creating diet individual restrictions", "error", err)
		return nil, err
	}

	return msg, nil
}

func (c *Commands) FoodConfigurationHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
А теперь перейдем к заполнению избранных продуктов.

❔ Начнем с круп для завтрака. Опишите какие крупы вы предпочитаете на завтрак

Например: манка,рисовая,пшенная,овсянная,кукурузная
`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventCreateFoodConfiguration); err != nil {
		c.logger.Errorw("Error creating diet food configuration", "error", err)
	}

	return msg, nil
}

func (c *Commands) CompleteConfigurationHandler(
	ctx context.Context,
	meta *entity.TelegramMeta,
	chat *entity.Chat,
	dietConfiguration *entity.DietConfiguration,
	botFSM *flow.BotFSM,
) (tgbotapi.Chattable, error) {
	text := `
	Ура! Конфигурация завершена!
	
	Теперь можно приступать к формированию рациона.
	`

	msg := tgbotapi.NewMessage(meta.ChatID, text)

	if err := botFSM.Event(flow.EventMainMenu); err != nil {
		c.logger.Errorw("Error creating diet food configuration", "error", err)
	}

	return msg, nil
}
