package commands

import "diet_bot/internal/entity"

type Repository interface {
	UpsertUser(user *entity.User) error
	GetUser(id int64) (*entity.User, error)

	GetChat(id int64) (*entity.Chat, error)
	UpsertChat(chat *entity.Chat) error

	UpsertDietConfiguration(dietConfiguration *entity.DietConfiguration) error
	GetDietConfiguration(userID int64) (*entity.DietConfiguration, error)
}
