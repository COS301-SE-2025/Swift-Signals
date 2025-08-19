package service

import (
	"errors"
	"strings"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func checkPassword(inputPassword, storedHashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(inputPassword))
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func handleValidationError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMap := make(map[string]string)
		for _, fieldError := range validationErrors {
			errorMap[strings.ToLower(fieldError.Field())] = getValidationErrorMessage(fieldError)
		}
		return errs.NewValidationError(
			"invalid input",
			map[string]any{"validation errors": errorMap},
		)
	}

	return errs.NewInternalError("validation failed", err, nil)
}

func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters long"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters long"
	case "uuid4":
		return fe.Field() + " must be a valid UUID"
	default:
		return fe.Field() + " is invalid"
	}
}
