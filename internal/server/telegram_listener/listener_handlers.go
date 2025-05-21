package telegramlistener

import (
	"context"
	generatediet "diet_bot/internal/commands/generate-diet"
	seediet "diet_bot/internal/commands/see-diet"
	"diet_bot/internal/flow"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (l *Listener) handleCommands(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	l.logger.With(
		"chat_id", update.Message.Chat.ID,
		"user_id", update.Message.From.ID,
		"message_id", update.Message.MessageID,
	).Info("Received new message")

	switch update.Message.Text {
	case flow.CommandStart:
		l.logger.Info("User send the start command")

		msg, err := l.commands.StartHandler(context.Background(), update)
		if err != nil {
			return nil, err
		}

		return msg, nil
	case flow.CommandMenu:
		l.logger.Info("User send the menu command")

		msg, err := l.commands.MenuHandler(context.Background(), update)
		if err != nil {
			return nil, err
		}

		return msg, nil
	case flow.CommandGenerateDiet:
		l.logger.Info("User send the generate diet command")

		msg := l.commands.GenerateDietHandler(context.Background(), update)

		return msg, nil
	default:
		if l.commands.IsFillUserConfigurationFlow(update) {
			msgs := l.commands.ExecuteFillment(context.Background(), update)

			for _, msg := range msgs {
				_, err := l.bot.Send(msg)
				if err != nil {
					l.logger.Error("Failed to send message", zap.Error(err))
				}
			}

			return nil, nil
		}

		l.logger.Infow("Unknown command received",
			"command", update.Message.Text,
			"user_id", update.Message.From.ID)

		return nil, nil
	}
}

func (l *Listener) handleCallback(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	l.bot.Request(callback)

	data := update.CallbackQuery.Data

	l.logger.With(
		"chat_id", update.CallbackQuery.Message.Chat.ID,
		"user_id", update.CallbackQuery.From.ID,
		"data", data,
	).Info("Received new callback")

	switch {
	case data == flow.CommandMenu:
		l.logger.Info("User pressed the menu button")

		deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
		_, err := l.bot.Send(deleteMsg)
		if err != nil {
			l.logger.Error("Failed to delete message", zap.Error(err))
		}

		msg, err := l.commands.MenuHandler(context.Background(), update)
		if err != nil {
			return nil, err
		}

		return msg, nil
	case data == flow.CommandSeeDietProducts:
		l.logger.Info("User pressed the see diet products button")

		l.deleteMessage(update.CallbackQuery.Message.MessageID, update)

		msg := l.commands.SeeDietProductsHandler(context.Background(), update)

		return msg, nil
	case data == flow.CommandGenerateDiet:
		l.logger.Info("User pressed the generate diet button")

		l.deleteMessage(update.CallbackQuery.Message.MessageID, update)

		msg := l.commands.GenerateDietHandler(context.Background(), update)

		return msg, nil
	case generatediet.CommandHasGenerateDietDays(data):
		l.logger.Info("User pressed the generate diet days button")

		l.deleteMessage(update.CallbackQuery.Message.MessageID, update)

		waitText := "Щас как наколдуем тебе красоту... ✨"
		waitMsg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, waitText)

		sentMsg, err := l.bot.Send(waitMsg)
		if err != nil {
			l.logger.Error("Failed to send wait message", zap.Error(err))
		}

		msg := l.commands.GenerateDietDaysHandler(
			context.Background(),
			update,
		)

		l.deleteMessageByMsgID(&sentMsg)

		return msg, nil
	case data == flow.CommandSeeDiet:
		l.logger.Info("User pressed the see diet button")

		l.deleteMessage(update.CallbackQuery.Message.MessageID, update)

		msg := l.commands.SeeDietHandler(context.Background(), update)

		return msg, nil
	case seediet.CommandHasDietDay(data):
		l.logger.Info("User pressed the see diet day button")

		l.deleteMessage(update.CallbackQuery.Message.MessageID, update)

		msg := l.commands.SeeDietDayHandler(context.Background(), update)

		return msg, nil
	case data == flow.CommandStartFillUserConfig:
		l.logger.Info("User pressed the start fill user config button")

		msg, err := l.commands.StartFlow(context.Background(), update)
		if err != nil {
			return nil, err
		}

		return msg, nil
	default:
		if l.commands.IsFillUserConfigurationFlow(update) {
			msgs := l.commands.ExecuteFillment(context.Background(), update)

			for _, msg := range msgs {
				_, err := l.bot.Send(msg)
				if err != nil {
					l.logger.Error("Failed to send message", zap.Error(err))
				}
			}

			return nil, nil
		}

		l.logger.Infow("Unknown callback data received",
			"data", update.CallbackQuery.Data,
			"user_id", update.CallbackQuery.From.ID)

		return nil, nil
	}
}

func (l *Listener) deleteMessage(_ int, update *tgbotapi.Update) {
	deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	_, err := l.bot.Send(deleteMsg)
	if err != nil {
		l.logger.Error("Failed to delete message", zap.Error(err))
	}
}

func (l *Listener) deleteMessageByMsgID(msg *tgbotapi.Message) {
	deleteMsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	_, err := l.bot.Send(deleteMsg)
	if err != nil {
		l.logger.Error("Failed to delete message", zap.Error(err))
	}
}

func (l *Listener) errorHandler(err error, chatID int64) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(chatID, "Ой, кажется что-то пошло не так. Попробуйте еще раз")

	l.bot.Send(msg)

	return msg
}
