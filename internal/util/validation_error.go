package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/grachmannico95/mileapp-test-be/internal/dto"
)

func ParseValidationError(err error) []dto.ErrorItem {
	var errors []dto.ErrorItem

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, dto.ErrorItem{
				Field:   getFieldName(fieldError.Field()),
				Message: getErrorMessage(fieldError),
			})
		}
		return errors
	}

	return []dto.ErrorItem{
		{
			Message: err.Error(),
		},
	}
}

func getFieldName(field string) string {
	return strings.ToLower(field)
}

func getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "Invalid email format"
	case "min":
		// Check if it's a string length validation or number validation
		switch fe.Kind().String() {
		case "string":
			return fmt.Sprintf("%s must be at least %s characters", field, param)
		default:
			return fmt.Sprintf("%s must be at least %s", field, param)
		}
	case "max":
		// Check if it's a string length validation or number validation
		switch fe.Kind().String() {
		case "string":
			return fmt.Sprintf("%s must not exceed %s characters", field, param)
		default:
			return fmt.Sprintf("%s must not exceed %s", field, param)
		}
	case "oneof":
		// Format the options nicely
		options := strings.ReplaceAll(param, " ", ", ")
		return fmt.Sprintf("%s must be one of: %s", field, options)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, param)
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be a numeric value", field)
	case "url":
		return "Invalid URL format"
	case "uri":
		return "Invalid URI format"
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime in format %s", field, param)
	default:
		return fmt.Sprintf("Invalid value for %s", getFieldName(field))
	}
}
