package validator

import (
	"errors"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validatable interface {
	ErrorMessages() map[string]string
}

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func Validate(c *fiber.Ctx, dto any) error {
	if err := c.BodyParser(dto); err != nil {
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "invalid character") || strings.Contains(errorMsg, "unexpected end") {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		return fiber.NewError(fiber.StatusBadRequest, "Failed to process your request")
	}

	validate := validator.New()

	if err := validate.Struct(dto); err != nil {
		errorBag := make(map[string]string)

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {

			customMessages := make(map[string]string)

			if v, ok := dto.(Validatable); ok {
				customMessages = v.ErrorMessages()
			}

			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				tag := fieldErr.Tag()

				jsonField := field
				t := reflect.TypeOf(dto)
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				if sf, ok := t.FieldByName(field); ok {
					jsonTag := sf.Tag.Get("json")
					if jsonTag != "" && jsonTag != "-" {
						jsonField = strings.Split(jsonTag, ",")[0]
					}
				}

				key := field + "." + tag

				if msg, exists := customMessages[key]; exists {
					errorBag[jsonField] = msg
				} else {
					errorBag[jsonField] = generateDefaultMessage(field, tag, fieldErr)
				}
			}
		}

		return &ValidationError{Errors: errorBag}
	}

	return nil
}

func generateDefaultMessage(field, tag string, fieldErr validator.FieldError) string {
	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " must have at least " + fieldErr.Param() + " items"
	case "max":
		return field + " cannot have more than " + fieldErr.Param() + " items"
	case "dive":
		return "Invalid items in " + field
	default:
		return field + " is invalid"
	}
}
