package entity

import (
	"time"

	"github.com/google/uuid"
)

type GeneratedDiet struct {
	ID        uuid.UUID   `json:"id" bson:"_id,omitempty"`
	UserID    string      `json:"user_id" bson:"user_id"`
	ConfigID  string      `json:"config_id" bson:"config_id"`
	DailyDiet []DailyDiet `json:"daily_diet" bson:"daily_diet"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}

type DailyDiet struct {
	ID            uuid.UUID `json:"id" bson:"_id,omitempty"`
	TotalCalories int       `json:"total_calories" bson:"total_calories"`
	TotalPFC      PFC       `json:"total_pfc" bson:"total_pfc"`
	Meals         []Meal    `json:"meals" bson:"meals"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

type Meal struct {
	ID        uuid.UUID `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Time      string    `json:"time" bson:"time"`
	Calories  int       `json:"calories" bson:"calories"`
	PFC       PFC       `json:"pfc" bson:"pfc"`
	Products  []Product `json:"products" bson:"products"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Product struct {
	ID        uuid.UUID `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Calories  int       `json:"calories" bson:"calories"`
	PFC       PFC       `json:"pfc" bson:"pfc"`
	Weight    int       `json:"weight" bson:"weight"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}