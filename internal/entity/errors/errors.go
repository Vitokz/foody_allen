package entity

import (
	internalerrors "diet_bot/internal/lib/errors"
	"errors"
)

var (
	ErrorUserNotFound = internalerrors.NewDatabaseError(errors.New("user not found"))

	ErrorFailedToSaveUser = internalerrors.NewDatabaseError(errors.New("failed to save user"))
	ErrorFailedToSaveChat = internalerrors.NewDatabaseError(errors.New("failed to save chat"))
	ErrorFailedToGetUser  = internalerrors.NewDatabaseError(errors.New("failed to get user"))

	ErrorInvalidAllergies = errors.New("invalid allergies")
	ErrorInvalidMealTypes = errors.New("invalid meal types")
)
