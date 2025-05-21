package userconfig

import (
	"context"
	"diet_bot/internal/entity"
	"diet_bot/internal/flow"
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FillConfig struct {
	CurrentState      string
	NextEvent         string
	PromptMessage     func(chatID int64) tgbotapi.Chattable
	Validation        func(chatID int64, value string) tgbotapi.Chattable
	FieldSetter       func(config *entity.UserConfiguration, value string) (response *string, err error)
	IsWaitingCallback bool
}

var fillConfig = []FillConfig{
	{
		CurrentState: flow.StateUserConfiguration_Height,
		NextEvent:    flow.EventUserConfigurationWeight,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			return tgbotapi.NewMessage(chatID, "✍️ Введи свой рост в сантиметрах")
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			height, err := strconv.Atoi(value)
			if err != nil {
				return tgbotapi.NewMessage(chatID, "Рост должен быть числом! Вот эти краказябры себе оставь")
			}

			if height < 50 || height > 250 {
				return tgbotapi.NewMessage(chatID, "Рост должен быть в диапазоне 50-250 см")
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			height, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("рост должен быть числом")
			}

			config.Height = height

			response := fmt.Sprintf("Рост успешно установлен: %d см", height)

			return &response, nil
		},
	},
	{
		CurrentState: flow.StateUserConfiguration_Weight,
		NextEvent:    flow.EventUserConfigurationGender,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			return tgbotapi.NewMessage(chatID, "✍️ Введи свой вес в килограммах")
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			value = strings.ReplaceAll(value, ",", ".")
			weight, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return tgbotapi.NewMessage(chatID, "Вес должен быть числом! Вот эти краказябры себе оставь")
			}

			if weight < 35 || weight > 300 {
				return tgbotapi.NewMessage(chatID, "Вес должен быть в диапазоне 35-300 кг")
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			value = strings.ReplaceAll(value, ",", ".")
			weight, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, errors.New("вес должен быть числом")
			}

			config.Weight = weight
			response := fmt.Sprintf("Вес успешно установлен: %.2f кг", weight)
			return &response, nil
		},
	},
	{
		CurrentState: flow.StateUserConfiguration_Gender,
		NextEvent:    flow.EventUserConfigurationAge,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			msg := tgbotapi.NewMessage(chatID, "👇 Выбери пол")
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("👨‍💼 Мужской", flow.CommandFillUserConfigGenderMale),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("👩‍💼 Женский", flow.CommandFillUserConfigGenderFemale),
				),
			)
			msg.ReplyMarkup = &replyMarkup

			return msg
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			if !isGenderFillCommand(value) {
				msg := tgbotapi.NewMessage(chatID, "👇 Выбери пожалуйста из списка")
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("👨‍💼 Мужской", flow.CommandFillUserConfigGenderMale),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("👩‍💼 Женский", flow.CommandFillUserConfigGenderFemale),
					),
				)
				msg.ReplyMarkup = &replyMarkup

				return msg
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			config.Gender = getGenderFromCommand(value)
			response := fmt.Sprintf("Пол успешно установлен: %s", entity.GenderToText(config.Gender))
			return &response, nil
		},
		IsWaitingCallback: true,
	},
	{
		CurrentState: flow.StateUserConfiguration_Age,
		NextEvent:    flow.EventUserConfigurationGoal,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			return tgbotapi.NewMessage(chatID, "✍️ Введи свой возраст")
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			age, err := strconv.Atoi(value)
			if err != nil {
				return tgbotapi.NewMessage(chatID, "Возраст должен быть целым числом! Ты где такие числа видел????")
			}

			if age < 12 || age > 100 {
				return tgbotapi.NewMessage(chatID, "Возраст должен быть в диапазоне 12-100 лет")
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			age, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("возраст должен быть числом")
			}

			config.Age = age
			response := fmt.Sprintf("Возраст успешно установлен: %d лет", age)
			return &response, nil
		},
	},
	{
		CurrentState: flow.StateUserConfiguration_Goal,
		NextEvent:    flow.EventUserConfigurationActivity,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			msg := tgbotapi.NewMessage(chatID, "👇 Какую цель преследуешь?")
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("⚖️ Похудеть", flow.CommandFillDietGoalLoseWeight),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("💪 Быть в форме", flow.CommandFillDietGoalMaintainWeight),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🏋️‍♂️ Набрать массу", flow.CommandFillDietGoalGainWeight),
				),
			)
			msg.ReplyMarkup = &replyMarkup

			return msg
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			if !isGoalFillCommand(value) {
				msg := tgbotapi.NewMessage(chatID, "👇 Выбери пожалуйста из списка")
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("⚖️ Похудеть", flow.CommandFillDietGoalLoseWeight),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("💪 Быть в форме", flow.CommandFillDietGoalMaintainWeight),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🏋️‍♂️ Набрать массу", flow.CommandFillDietGoalGainWeight),
					),
				)
				msg.ReplyMarkup = &replyMarkup

				return msg
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			config.Goal = getGoalFromCommand(value)
			response := fmt.Sprintf("Цель успешно установлена: %s", entity.GoalToText(config.Goal))

			return &response, nil
		},
		IsWaitingCallback: true,
	},
	{
		CurrentState: flow.StateUserConfiguration_Activity,
		NextEvent:    flow.EventUserConfigurationDietType,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			msg := tgbotapi.NewMessage(chatID, `
👇 Выбери свой уровень активности:

🪑 Сидячий — мало двигаюсь, работаю в офисе.
🏃‍♂️ Малоактивный — 1-3 раза в неделю занимаюсь спортом.
🏋️‍♂️ Умеренно активный — 3-5 раз в неделю занимаюсь спортом.
💪 Активный — 6-7 раз в неделю занимаюсь спортом.
🦾 Очень активный — тренировки и физическая работа.
		`)
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🪑 Сидячий", flow.CommandFillUserConfigActivitySedentary),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🏃‍♂️ Малоактивный", flow.CommandFillUserConfigActivityLow),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🏋️‍♂️ Умеренно активный", flow.CommandFillUserConfigActivityMedium),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("💪 Активный", flow.CommandFillUserConfigActivityHigh),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🦾 Очень активный", flow.CommandFillUserConfigActivityVeryHigh),
				),
			)
			msg.ReplyMarkup = &replyMarkup

			return msg
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			if !isActivityFillCommand(value) {
				msg := tgbotapi.NewMessage(chatID, "👇 Выбери пожалуйста из списка")
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🪑 Сидячий", flow.CommandFillUserConfigActivitySedentary),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🏃‍♂️ Малоактивный", flow.CommandFillUserConfigActivityLow),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🏋️‍♂️ Умеренно активный", flow.CommandFillUserConfigActivityMedium),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("💪 Активный", flow.CommandFillUserConfigActivityHigh),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🦾 Очень активный", flow.CommandFillUserConfigActivityVeryHigh),
					),
				)
				msg.ReplyMarkup = &replyMarkup

				return msg
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			config.Activity = getActivityFromCommand(value)

			response := fmt.Sprintf("Уровень активности: %s", entity.ActivityToText(config.Activity))

			return &response, nil
		},
		IsWaitingCallback: true,
	},
	{
		CurrentState: flow.StateUserConfiguration_DietType,
		NextEvent:    flow.EventUserConfigurationAllergies,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			msg := tgbotapi.NewMessage(chatID, `
Может быть ты следуешь какой-то особенной диете? 

👇 Выбери из списка:
`)
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Нет, не следую", flow.CommandFillUserConfigDietTypeAnything),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Кето", flow.CommandFillUserConfigDietTypeKeto),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Палео", flow.CommandFillUserConfigDietTypePaleo),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Веган", flow.CommandFillUserConfigDietTypeVegan),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Вегетарианская", flow.CommandFillUserConfigDietTypeVegetarian),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Средиземноморская", flow.CommandFillUserConfigDietTypeMediterranean),
				),
			)
			msg.ReplyMarkup = &replyMarkup

			return msg
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			if !isDietTypeFillCommand(value) {
				msg := tgbotapi.NewMessage(chatID, "👇 Выбери пожалуйста из списка")
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Нет, не следую", flow.CommandFillUserConfigDietTypeAnything),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Кето", flow.CommandFillUserConfigDietTypeKeto),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Палео", flow.CommandFillUserConfigDietTypePaleo),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Веган", flow.CommandFillUserConfigDietTypeVegan),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Вегетарианская", flow.CommandFillUserConfigDietTypeVegetarian),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Средиземноморская", flow.CommandFillUserConfigDietTypeMediterranean),
					),
				)
				msg.ReplyMarkup = &replyMarkup

				return msg
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			config.DietType = getDietTypeFromCommand(value)
			response := fmt.Sprintf("Тип диеты: %s", entity.DietTypeToText(config.DietType))
			return &response, nil
		},
		IsWaitingCallback: true,
	},
	{
		CurrentState: flow.StateUserConfiguration_Allergies,
		NextEvent:    flow.EventUserConfigurationMealTypes,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			return tgbotapi.NewMessage(chatID, `
✍️ Есть ли у тебя пищевая аллергия?
			
	1. Глютен
	2. Арахис
	3. Яйца
	4. Рыба
	5. Орехи
	6. Молочка
	7. Соя
	8. Ракообразные
			
Если есть, то напиши их через запятую. Например:

"глютен, арахис, яйца" или просто "глютен"

Если ничего такого нету, пиши "нет".
		`)
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			if value == "нет" {
				return nil
			}

			_, err := entity.AllergensFromTextToEntity(value)
			if err != nil {
				return tgbotapi.NewMessage(chatID, "Пожалуйста, введи аллергии через запятую. Например: 'глютен,арахис,яйца'")
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			response := fmt.Sprintf("Аллергии успешно установлены: %s", value)

			if value == "нет" {
				config.Allergies = nil
				return &response, nil
			}

			allergies, err := entity.AllergensFromTextToEntity(value)
			if err != nil {
				return nil, err
			}

			config.Allergies = allergies

			return &response, nil
		},
	},
	{
		CurrentState: flow.StateUserConfiguration_MealTypes,
		NextEvent:    flow.EventMainMenu,
		PromptMessage: func(chatID int64) tgbotapi.Chattable {
			return tgbotapi.NewMessage(chatID, `
✍️ Введи приемы пищи. 
			
	Варианты:
	1. Завтрак
	2. Обед
	3. Ужин
	4. Перекус

Например:

"завтрак, обед, ужин"
`)
		},
		Validation: func(chatID int64, value string) tgbotapi.Chattable {
			meals, err := entity.MealTypesFromTextToEntity(value)
			if err != nil {
				return tgbotapi.NewMessage(chatID, "Пожалуйста, введи приемы пищи через запятую. Например: 'завтрак,обед,ужин'")
			}

			if len(meals) == 0 {
				return tgbotapi.NewMessage(chatID, "Пожалуйста, введи приемы пищи через запятую. Например: 'завтрак,обед,ужин'")
			}

			return nil
		},
		FieldSetter: func(config *entity.UserConfiguration, value string) (*string, error) {
			meals, err := entity.MealTypesFromTextToEntity(value)
			if err != nil {
				return nil, err
			}

			config.MealTypes = meals

			response := fmt.Sprintf("Приемы пищи успешно установлены: %s", value)
			return &response, nil
		},
	},
}

var flowConfigMap map[string]FillConfig

func init() {
	flowConfigMap = make(map[string]FillConfig)
	for _, config := range fillConfig {
		flowConfigMap[config.CurrentState] = config
	}
}

func (c *Commands) IsFillUserConfigurationFlow(update *tgbotapi.Update) bool {
	meta := entity.NewMeta(update)

	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("Error getting chat", "error", err)
		return false
	}

	ok := flow.IsCreateUserConfigurationFillState(chat.State)

	if !ok {
		c.logger.Infof("Chat state: %s is not fill user configuration", chat.State)
	} else {
		c.logger.Infof("Chat state: %s is fill user configuration", chat.State)
	}

	return ok
}

func (c *Commands) ExecuteFillment(
	ctx context.Context,
	update *tgbotapi.Update,
) []tgbotapi.Chattable {
	meta := entity.NewMeta(update)

	// 1. Get chat
	chat, err := c.repository.GetChat(meta.ChatID)
	if err != nil {
		c.logger.Errorw("error getting chat", "error", err)
		return []tgbotapi.Chattable{
			tgbotapi.NewMessage(meta.ChatID, "произошла ошибка, нажмите /start еще раз"),
		}
	}

	// 2. Get user config
	userConfig, err := c.repository.GetUserConfiguration(chat.UserID)
	if err != nil {
		c.logger.Errorw("error getting user configuration", "error", err)
		return []tgbotapi.Chattable{
			tgbotapi.NewMessage(meta.ChatID, "произошла ошибка, нажмите /start еще раз"),
		}
	}

	// 3. Get config for current state
	config, exists := flowConfigMap[chat.State]
	if !exists {
		return []tgbotapi.Chattable{
			tgbotapi.NewMessage(meta.ChatID, "произошла ошибка при формировании конфигурации"),
		}
	}

	// 4. Validate input
	if config.IsWaitingCallback {
		if meta.CallbackData == nil {
			c.logger.Error("error expected not nil callback data")

			return []tgbotapi.Chattable{
				config.Validation(meta.ChatID, meta.Message.Text),
			}
		}

		if errMsg := config.Validation(meta.ChatID, *meta.CallbackData); errMsg != nil {
			return []tgbotapi.Chattable{
				errMsg,
			}
		}
	} else {
		if errMsg := config.Validation(meta.ChatID, meta.Message.Text); errMsg != nil {
			return []tgbotapi.Chattable{
				errMsg,
			}
		}
	}

	// 5. Set field value
	var userChoice string
	if config.IsWaitingCallback {
		if meta.CallbackData == nil {
			c.logger.Error("error expected not nil callback data")
			return []tgbotapi.Chattable{
				tgbotapi.NewMessage(meta.ChatID, "произошла ошибка при получении данных"),
			}
		}
		userChoice = *meta.CallbackData
		response, err := config.FieldSetter(userConfig, userChoice)
		if err != nil {
			return []tgbotapi.Chattable{
				tgbotapi.NewMessage(meta.ChatID, err.Error()),
			}
		}
		userChoice = *response
	} else {
		userChoice = meta.Message.Text
		response, err := config.FieldSetter(userConfig, userChoice)
		if err != nil {
			return []tgbotapi.Chattable{
				tgbotapi.NewMessage(meta.ChatID, err.Error()),
			}
		}
		userChoice = *response
	}

	// 6. Save changes
	if err := c.repository.SaveUserConfiguration(userConfig); err != nil {
		c.logger.Errorw("error saving user configuration", "error", err)
		return []tgbotapi.Chattable{
			tgbotapi.NewMessage(meta.ChatID, "произошла ошибка при сохранении данных"),
		}
	}

	messages := []tgbotapi.Chattable{}

	// Edit previous message with user's choice
	if !config.IsWaitingCallback {
		deleteMsg := tgbotapi.NewDeleteMessage(meta.ChatID, meta.Message.MessageID)
		editMsg := tgbotapi.NewEditMessageText(meta.ChatID, meta.Message.MessageID-1, fmt.Sprintf("✅ %s", userChoice))

		messages = append(messages, deleteMsg, editMsg)
	} else {
		editMsg := tgbotapi.NewEditMessageText(meta.ChatID, meta.Message.MessageID, fmt.Sprintf("✅ %s", userChoice))
		messages = append(messages, editMsg)
	}

	// 7. Transition to next state
	botFSM := flow.NewBotFSM(chat)
	if config.NextEvent == flow.EventMainMenu {
		if err := botFSM.Event(flow.EventMainMenu); err != nil {
			c.logger.Errorw("error transitioning to main menu", "error", err)
			return []tgbotapi.Chattable{
				tgbotapi.NewMessage(meta.ChatID, "произошла ошибка при переходе к главному меню"),
			}
		}

		if err := c.repository.UpsertChat(chat); err != nil {
			c.logger.Errorw("error saving chat", "error", err)
			return []tgbotapi.Chattable{
				tgbotapi.NewMessage(meta.ChatID, "произошла ошибка при сохранении данных"),
			}
		}

		msg := tgbotapi.NewMessage(meta.ChatID, "Конфигурация завершена!")
		messages = append(messages, msg)

		return messages
	}

	if err := botFSM.Event(config.NextEvent); err != nil {
		c.logger.Errorw("error transitioning to next state", "error", err)
		return []tgbotapi.Chattable{
			tgbotapi.NewMessage(meta.ChatID, "произошла ошибка при переходе к следующему шагу"),
		}
	}

	// 8. Save chat state
	if err := c.repository.UpsertChat(chat); err != nil {
		c.logger.Errorw("error saving chat", "error", err)
		return []tgbotapi.Chattable{
			tgbotapi.NewMessage(meta.ChatID, "произошла ошибка при сохранении данных"),
		}
	}

	// Return prompt for next state
	nextConfig := flowConfigMap[chat.State]

	messages = append(messages, nextConfig.PromptMessage(meta.ChatID))

	return messages
}

func isGenderFillCommand(command string) bool {
	return strings.HasPrefix(command, flow.CommandFillUserConfigGender)
}

func getGenderFromCommand(command string) string {
	return strings.TrimPrefix(command, flow.CommandFillUserConfigGender)
}

func isGoalFillCommand(command string) bool {
	return strings.HasPrefix(command, flow.CommandFillDietGoal)
}

func getGoalFromCommand(command string) string {
	return strings.TrimPrefix(command, flow.CommandFillDietGoal)
}

func isActivityFillCommand(command string) bool {
	return strings.HasPrefix(command, flow.CommandFillUserConfigActivity)
}

func getActivityFromCommand(command string) string {
	return strings.TrimPrefix(command, flow.CommandFillUserConfigActivity)
}

func isDietTypeFillCommand(command string) bool {
	return strings.HasPrefix(command, flow.CommandFillUserConfigDietType)
}

func getDietTypeFromCommand(command string) string {
	return strings.TrimPrefix(command, flow.CommandFillUserConfigDietType)
}
