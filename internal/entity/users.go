package entity

import "time"

type User struct {
	ID           int64     `json:"id" bson:"_id"`
	FirstName    string    `json:"first_name" bson:"first_name"`
	LastName     string    `json:"last_name" bson:"last_name"`
	Username     string    `json:"username" bson:"username"`
	LanguageCode string    `json:"language_code" bson:"language_code"`
	IsBot        bool      `json:"is_bot" bson:"is_bot"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

func (u *User) CollectionName() string {
	return "users"
}
