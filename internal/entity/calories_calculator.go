package entity

import "github.com/google/uuid"

type UserCalories struct {
	ID                  uuid.UUID `json:"id" bson:"_id"`
	UserID              int64     `json:"user_id" bson:"user_id"`
	BaseCalories        int64     `json:"base_calories" bson:"base_calories"`
	GoalCalories        int64     `json:"goal_calories" bson:"goal_calories"`
	GoalCaloriesPercent int64     `json:"goal_calories_percent" bson:"goal_calories_percent"`
	// true - увеличить, false - уменьшить
	GoalCaloriesDirection bool `json:"goal_calories_direction" bson:"goal_calories_direction"`
}

func (u *UserCalories) CollectionName() string {
	return "user_calories"
}

// calculateCaloriesMifflinStJeor рассчитывает базовый и дневной расход калорий по формуле Миффлина-Сан Жеора
func calculateCaloriesMifflinStJeor(cfg *UserConfiguration) (total float64) {
	activityMap := map[string]float64{
		ActivitySedentary:  1.2,
		ActivityLight:      1.375,
		ActivityModerate:   1.55,
		ActivityActive:     1.725,
		ActivityVeryActive: 1.9,
	}

	activityCoeff, ok := activityMap[cfg.Activity]
	if !ok {
		activityCoeff = 1.2 // по умолчанию малоподвижный
	}

	var bmr float64
	if cfg.Gender == GenderMale {
		bmr = 10*cfg.Weight + 6.25*float64(cfg.Height) - 5*float64(cfg.Age) + 5
	} else {
		bmr = 10*cfg.Weight + 6.25*float64(cfg.Height) - 5*float64(cfg.Age) - 161
	}
	total = bmr * activityCoeff
	return
}

func (u *UserCalories) CalculateCalories(userCfg *UserConfiguration) {
	baseCalories := calculateCaloriesMifflinStJeor(userCfg)

	var goalCalories int64
	var goalCaloriesDirection bool
	switch userCfg.Goal {
	case GoalLoseWeight:
		goalCalories = int64(baseCalories * 0.9)
		goalCaloriesDirection = false
	case GoalGainWeight:
		goalCalories = int64(baseCalories)
		goalCaloriesDirection = true
	case GoalMaintainWeight:
		goalCalories = int64(baseCalories * 1.1)
		goalCaloriesDirection = true
	}

	u.ID = uuid.New()
	u.UserID = userCfg.UserID
	u.BaseCalories = int64(baseCalories)
	u.GoalCalories = goalCalories
	u.GoalCaloriesDirection = goalCaloriesDirection
}
