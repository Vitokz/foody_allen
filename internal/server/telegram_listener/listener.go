package telegramlistener

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"diet_bot/internal/commands"
	generatediet "diet_bot/internal/commands/generate-diet"
	seediet "diet_bot/internal/commands/see-diet"
	"diet_bot/internal/flow"
)

type Listener struct {
	bot      *tgbotapi.BotAPI
	commands *commands.Commands

	logger *zap.SugaredLogger

	exitChan chan struct{}
}

func NewListener(
	botToken string,
	commands *commands.Commands,
	logger *zap.SugaredLogger,
) (*Listener, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	setupBotCommands(bot)

	return &Listener{
		bot:      bot,
		commands: commands,
		logger:   logger,
		exitChan: make(chan struct{}),
	}, nil
}

func (l *Listener) Listen() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := l.bot.GetUpdatesChan(u)

	for {
		select {
		case <-l.exitChan:
			return nil
		case update := <-updates:
			if update.Message != nil {
				switch update.Message.Text {
				case flow.CommandStart:
					l.logger.Info("User started the bot")

					msg := l.commands.StartHandler(context.Background(), &update)

					l.bot.Send(msg)
				case flow.CommandMenu:
					l.logger.Info("User pressed the menu button")

					msg := l.commands.MenuHandler(context.Background(), &update)

					l.bot.Send(msg)
				case flow.CommandGenerateDiet:
					l.logger.Info("User pressed the generate diet button")

					msg := l.commands.GenerateDietHandler(context.Background(), &update)

					l.bot.Send(msg)
				default:
					if l.commands.IsFillDiet(context.Background(), &update) {
						msg := l.commands.FillDiet(context.Background(), &update)

						l.bot.Send(msg)

					} else if l.commands.IsFillProducts(context.Background(), &update) {
						msg := l.commands.FillDiet(context.Background(), &update)

						l.bot.Send(msg)
					}
				}
			} else if update.CallbackQuery != nil {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
				l.bot.Request(callback)

				data := update.CallbackQuery.Data

				switch {
				case data == flow.CommandMenu:
					l.logger.Info("User pressed the menu button")

					// 1. Удаляем сообщение, в котором была нажата кнопка
					deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					_, err := l.bot.Send(deleteMsg)
					if err != nil {
						l.logger.Error("Failed to delete message", zap.Error(err))
					}

					// 2. Отправляем новое сообщение
					msg := l.commands.MenuHandler(context.Background(), &update)
					l.bot.Send(msg)
				case data == flow.CommandSeeDietProducts:
					l.logger.Info("User pressed the see diet products button")

					l.deleteMessage(update.CallbackQuery.Message.MessageID, &update)

					msg := l.commands.SeeDietProductsHandler(context.Background(), &update)

					l.bot.Send(msg)
				case data == flow.CommandGenerateDiet:
					l.logger.Info("User pressed the generate diet button")

					l.deleteMessage(update.CallbackQuery.Message.MessageID, &update)

					msg := l.commands.GenerateDietHandler(context.Background(), &update)

					l.bot.Send(msg)
				case generatediet.CommandHasGenerateDietDays(data):
					l.logger.Info("User pressed the generate diet days button")

					l.deleteMessage(update.CallbackQuery.Message.MessageID, &update)

					msg := l.commands.GenerateDietDaysHandler(
						context.Background(),
						&update,
					)

					l.bot.Send(msg)
				case data == flow.CommandFillConfig:
					l.logger.Info("User pressed the fill config button")

					l.deleteMessage(update.CallbackQuery.Message.MessageID, &update)

					msg := l.commands.CreateDietHandler(context.Background(), &update)

					l.bot.Send(msg)
				case data == flow.CommandSeeDiet:
					l.logger.Info("User pressed the see diet button")

					l.deleteMessage(update.CallbackQuery.Message.MessageID, &update)

					msg := l.commands.SeeDietHandler(context.Background(), &update)

					l.bot.Send(msg)
				case seediet.CommandHasDietDay(data):
					l.logger.Info("User pressed the see diet day button")

					l.deleteMessage(update.CallbackQuery.Message.MessageID, &update)

					msg := l.commands.SeeDietDayHandler(context.Background(), &update)

					l.bot.Send(msg)
				default:
					l.logger.Infow("Unknown callback data received",
						"data", update.CallbackQuery.Data,
						"user_id", update.CallbackQuery.From.ID)
				}
			}
		}
	}
}

func (l *Listener) Stop() {
	close(l.exitChan)
}

func (l *Listener) deleteMessage(_ int, update *tgbotapi.Update) {
	deleteMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	_, err := l.bot.Send(deleteMsg)
	if err != nil {
		l.logger.Error("Failed to delete message", zap.Error(err))
	}
}
