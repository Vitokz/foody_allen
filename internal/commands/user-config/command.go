package userconfig

import (
	"diet_bot/internal/repository"

	"go.uber.org/zap"
)

type Commands struct {
	repository *repository.Client
	logger     *zap.SugaredLogger
}

func NewCommands(repository *repository.Client, logger *zap.SugaredLogger) *Commands {
	return &Commands{
		repository: repository,
		logger:     logger,
	}
}
