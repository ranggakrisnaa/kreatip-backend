package config

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kreatip/kreatip-backend/internal/model"
)

type errorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type errorResponse struct {
	Errors []errorItem `json:"errors"`
}

// ErrorHandler is the global Fiber error handler.
// Priority: AppError → validator.ValidationErrors → fiber.Error → 500
func ErrorHandler(c *fiber.Ctx, err error) error {
	// 1. Business errors (from usecase layer)
	var appErr *model.AppError
	if errors.As(err, &appErr) {
		return c.Status(appErr.Code).JSON(errorResponse{
			Errors: []errorItem{{Field: appErr.Field, Message: appErr.Message}},
		})
	}

	// 2. Validation errors (from validator.Struct)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		items := make([]errorItem, len(ve))
		for i, fe := range ve {
			items[i] = errorItem{
				Field:   toSnakeCase(fe.Field()),
				Message: validationMessage(fe),
			}
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errorResponse{Errors: items})
	}

	// 3. Fiber built-in errors (404, 405, etc.)
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return c.Status(fe.Code).JSON(errorResponse{
			Errors: []errorItem{{Message: fe.Message}},
		})
	}

	// 4. Fallback — 500, hide internal detail
	return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{
		Errors: []errorItem{{Message: "internal server error"}},
	})
}

func validationMessage(fe validator.FieldError) string {
	field := toSnakeCase(fe.Field())
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " minimum length is " + fe.Param()
	case "max":
		return field + " maximum length is " + fe.Param()
	case "oneof":
		return field + " must be one of: " + fe.Param()
	case "uuid":
		return field + " must be a valid UUID"
	case "len":
		return field + " must be exactly " + fe.Param() + " characters"
	default:
		return field + " is invalid"
	}
}

// toSnakeCase converts "FieldName" → "field_name" for consistent field names in errors.
func toSnakeCase(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, c+32)
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}
