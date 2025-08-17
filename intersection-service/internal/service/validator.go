package service

import (
	"errors"
	"strings"

	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/go-playground/validator/v10"
)

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
