package config

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorHandler is the global Fiber error handler.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Validation errors
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errs := make([]ValidationError, len(ve))
		for i, fe := range ve {
			errs[i] = ValidationError{
				Field:   fe.Field(),
				Message: validationMessage(fe),
			}
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": errs,
		})
	}

	// Fiber errors
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return c.Status(fe.Code).JSON(fiber.Map{
			"errors": []ValidationError{{Message: fe.Message}},
		})
	}

	// Default 500
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"errors": []ValidationError{{Message: "internal server error"}},
	})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email"
	case "min":
		return fe.Field() + " minimum value is " + fe.Param()
	case "max":
		return fe.Field() + " maximum value is " + fe.Param()
	default:
		return fe.Field() + " is invalid"
	}
}
