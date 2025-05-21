package telegramlistener

import (
	"errors"
	"runtime/debug"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"diet_bot/internal/commands"
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

	var wg sync.WaitGroup

	// Ограничиваем количество одновременных обработчиков
	semaphore := make(chan struct{}, 10) // максимум 10 одновременных обработок

	for {
		select {
		case <-l.exitChan:
			// Дожидаемся завершения всех горутин перед выходом
			wg.Wait()
			return nil
		case update := <-updates:
			// Добавляем в WaitGroup новую задачу
			wg.Add(1)

			// Блокируем семафор
			semaphore <- struct{}{}

			// Обрабатываем каждое сообщение асинхронно
			go func(update tgbotapi.Update) {
				defer wg.Done()
				defer func() { <-semaphore }() // освобождаем семафор в любом случае

				// Обрабатываем паники
				defer func() {
					if r := recover(); r != nil {
						l.logger.Error("Panic in message handler",
							zap.Any("recover", r),
							zap.String("stack", string(debug.Stack())))

						l.errorHandler(errors.New("panic at during message handler"), update.Message.Chat.ID)
					}
				}()

				if update.Message != nil {
					msg, err := l.handleCommands(&update)
					if err != nil {
						l.logger.Error("Failed to handle command", zap.Error(err))
						l.errorHandler(err, update.Message.Chat.ID)

						return
					}

					if msg == nil {
						return
					}

					_, err = l.bot.Send(msg)
					if err != nil {
						l.logger.Error("Failed to send message", zap.Error(err))
					}
				} else if update.CallbackQuery != nil {
					msg, err := l.handleCallback(&update)
					if err != nil {
						l.logger.Error("Failed to handle callback", zap.Error(err))
						l.errorHandler(err, update.CallbackQuery.Message.Chat.ID)

						return
					}

					if msg == nil {
						return
					}

					_, err = l.bot.Send(msg)
					if err != nil {
						l.logger.Error("Failed to send message", zap.Error(err))
						l.errorHandler(err, update.Message.Chat.ID)

						return
					}
				}
			}(update)
		}
	}
}

func (l *Listener) Stop() {
	close(l.exitChan)
}
