package validation

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/utils"
)

func InitValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("Failed to initialize the validator engine")
	}

	RegisterCustomValidation(v)
	return nil
}

func HandleValidationErrors(err error) gin.H {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errors = make(map[string]string)

		for _, e := range validationErrors {
			root := strings.SplitN(e.Namespace(), ".", 2)[0]
			rawPath := strings.TrimPrefix(e.Namespace(), root+".")

			parts := strings.Split(rawPath, ".")
			for i := range parts {
				parts[i] = normalizeField(parts[i])
			}

			fieldPath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than %s", fieldPath, e.Param())
			case "gte":
				errors[fieldPath] = fmt.Sprintf("%s must be greater than or equal %s", fieldPath, e.Param())
			case "lt":
				errors[fieldPath] = fmt.Sprintf("%s must be less than %s", fieldPath, e.Param())
			case "lte":
				errors[fieldPath] = fmt.Sprintf("%s must be less than or equal %s", fieldPath, e.Param())
			case "uuid":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid UUID", fieldPath)
			case "required":
				errors[fieldPath] = fmt.Sprintf("%s is required", fieldPath)
			case "email":
				errors[fieldPath] = fmt.Sprintf("%s must be a valid email", fieldPath)
			case "datetime":
				errors[fieldPath] = fmt.Sprintf("%s must be in YYYY-MM-DD format", fieldPath)
			case "search":
				errors[fieldPath] = fmt.Sprintf("%s must contain only letters, numbers and spaces", fieldPath)
			case "slug":
				errors[fieldPath] = fmt.Sprintf("%s must contain only lowercase letters, numbers, hyphens, and dots", fieldPath)
			case "file_ext":
				allowedExts := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s only allow files with the specified extension: %s", fieldPath, allowedExts)
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[fieldPath] = fmt.Sprintf("%s must be one of the following values: %s", fieldPath, allowedValues)
			case "email_advanced":
				errors[fieldPath] = fmt.Sprintf("%s is on the banned list", fieldPath)
			case "password_strong":
				errors[fieldPath] = fmt.Sprintf("%s must be at least 8 characters long and include a lowercase letter, an uppercase letter, a number, and a special character.", fieldPath)
			}
		}

		return gin.H{"error": errors}
	}
	return gin.H{"error": "Invalid request: " + err.Error()}
}

func normalizeField(part string) string {
	base, index, found := strings.Cut(part, "[")
	if !found {
		return utils.CamelToSnake(part)
	}
	return utils.CamelToSnake(base) + "[" + index
}
