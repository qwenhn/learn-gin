package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func normalizeField(part string) string {
	base, index, found := strings.Cut(part, "[")
	if !found {
		return camelToSnake(part)
	}
	return camelToSnake(base) + "[" + index
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
			}
		}

		return gin.H{"error": errors}
	}
	return gin.H{"error": "Invalid request: " + err.Error()}
}

func ValidationInList(fieldName, value string, allowed map[string]bool) error {

	if !allowed[value] {
		return fmt.Errorf("%s mus be one of: %v", fieldName, getMapKeys(allowed))
	}
	return nil
}

func getMapKeys(m map[string]bool) []string {
	// Pre-allocate slice for performance
	keys := make([]string, 0, len(m))

	// Only capture the key by omitting the value variable
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("Failed to initialize the validator engine")
	}

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		filename := fl.Field().String()
		allowedFiles := fl.Param()

		if allowedFiles == "" {
			return false
		}

		allowedExts := strings.Fields(allowedFiles)
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")

		for _, allowed := range allowedExts {
			if ext == strings.ToLower(allowed) {
				return true
			}
		}

		return false
	})

	return nil
}
