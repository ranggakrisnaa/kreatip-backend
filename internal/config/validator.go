package config

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate {
	v := validator.New()
	// Register custom validators here as needed
	return v
}
