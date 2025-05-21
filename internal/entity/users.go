package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"

	internalerrors "diet_bot/internal/entity/errors"
)

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

const (
	GenderMale   = "male"
	GenderFemale = "female"

	GoalLoseWeight     = "lose_weight"
	GoalMaintainWeight = "maintain_weight"
	GoalGainWeight     = "gain_weight"

	ActivitySedentary  = "sedentary"
	ActivityLight      = "light"
	ActivityModerate   = "moderate"
	ActivityActive     = "active"
	ActivityVeryActive = "very_active"

	DietTypeAnything      = "anything"
	DietTypeKeto          = "keto"
	DietTypePaleo         = "paleo"
	DietTypeVegan         = "vegan"
	DietTypeVegetarian    = "vegetarian"
	DietTypeMediterranean = "mediterranean"

	AllergenGluten    = "gluten"
	AllergenPeanuts   = "peanuts"
	AllergenEggs      = "eggs"
	AllergenFish      = "fish"
	AllergenTreeNuts  = "tree_nuts"
	AllergenDairy     = "dairy"
	AllergenSoy       = "soy"
	AllergenShellfish = "shellfish"

	MealTypeBreakfast = "breakfast"
	MealTypeLunch     = "lunch"
	MealTypeDinner    = "dinner"
	MealTypeSnack     = "snack"
)

type UserConfiguration struct {
	ID        uuid.UUID `json:"id" bson:"_id"`
	UserID    int64     `json:"user_id" bson:"user_id"`
	Height    int       `json:"height" bson:"height"`
	Weight    float64   `json:"weight" bson:"weight"`
	Gender    string    `json:"gender" bson:"gender"` // male, female, other
	Age       int       `json:"age" bson:"age"`
	Goal      string    `json:"goal" bson:"goal"`
	Activity  string    `json:"activity" bson:"activity"`
	DietType  string    `json:"diet_type" bson:"diet_type"`
	Allergies []string  `json:"allergies" bson:"allergies"`   // gluten, peanuts, eggs, fish, tree_nuts, dairy, soy, shellfish,
	MealTypes []string  `json:"meal_types" bson:"meal_types"` // breakfast, lunch, dinner, snack
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	// MealTypePreferences map[string]string `json:"meal_type_preferences" bson:"meal_type_preferences"`
	// BodyFat             string            `json:"body_fat" bson:"body_fat"`
}

func (u *UserConfiguration) CollectionName() string {
	return "user_configurations"
}

func GoalToText(goal string) string {
	switch goal {
	case GoalLoseWeight:
		return "⚖️ Похудеть"
	case GoalMaintainWeight:
		return "💪 Быть в форме"
	case GoalGainWeight:
		return "🏋️‍♂️ Набрать массу"
	default:
		return "💪 Быть в форме"
	}
}

func ActivityToText(activity string) string {
	switch activity {
	case ActivitySedentary:
		return "🪑 Сидячий"
	case ActivityLight:
		return "🏃‍♂️ Малоактивный"
	case ActivityModerate:
		return "🏋️‍♂️ Умеренно активный"
	case ActivityActive:
		return "💪 Активный"
	case ActivityVeryActive:
		return "🦾 Очень активный"
	default:
		return "🪑 Сидячий"
	}
}

func DietTypeToText(dietType string) string {
	switch dietType {
	case DietTypeAnything:
		return "Любой"
	case DietTypeKeto:
		return "Кето"
	case DietTypePaleo:
		return "Пальео"
	case DietTypeVegan:
		return "Веган"
	case DietTypeVegetarian:
		return "Вегетариан"
	case DietTypeMediterranean:
		return "Средиземноморская"
	default:
		return "Любой"
	}
}

func GenderToText(gender string) string {
	switch gender {
	case GenderMale:
		return "👨‍💼 Мужской"
	case GenderFemale:
		return "👩‍💼 Женский"
	default:
		return "👨‍💼 Мужской"
	}
}

func AllergensFromTextToEntity(text string) ([]string, error) {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ToLower(text)

	rawAllergies := strings.Split(text, ",")
	allergies := make([]string, 0, len(rawAllergies))

	for _, allergy := range rawAllergies {
		switch strings.ToLower(allergy) {
		case "глютен":
			allergies = append(allergies, AllergenGluten)
		case "арахис":
			allergies = append(allergies, AllergenPeanuts)
		case "яйца":
			allergies = append(allergies, AllergenEggs)
		case "рыба":
			allergies = append(allergies, AllergenFish)
		case "орехи":
			allergies = append(allergies, AllergenTreeNuts)
		case "молочка":
			allergies = append(allergies, AllergenDairy)
		case "соя":
			allergies = append(allergies, AllergenSoy)
		case "ракообразные":
			allergies = append(allergies, AllergenShellfish)
		default:
			return nil, internalerrors.ErrorInvalidAllergies
		}
	}

	return allergies, nil
}

func MealTypesFromTextToEntity(text string) ([]string, error) {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ToLower(text)

	rawMealTypes := strings.Split(text, ",")
	mealTypes := make([]string, 0, len(rawMealTypes))

	for _, mealType := range rawMealTypes {
		switch strings.ToLower(mealType) {
		case "завтрак":
			mealTypes = append(mealTypes, MealTypeBreakfast)
		case "обед":
			mealTypes = append(mealTypes, MealTypeLunch)
		case "ужин":
			mealTypes = append(mealTypes, MealTypeDinner)
		case "перекус":
			mealTypes = append(mealTypes, MealTypeSnack)
		default:
			return nil, internalerrors.ErrorInvalidMealTypes
		}
	}

	return mealTypes, nil
}
