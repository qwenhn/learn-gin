package utils

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

func ValidationPositiveInt(fieldName string, value string) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number", fieldName)
	}

	if id <= 0 {
		return 0, fmt.Errorf("%s must be positive", fieldName)
	}

	return id, nil
}

func ValidationUuid(fieldName string, value string) (uuid.UUID, error) {
	uid, err := uuid.Parse(value)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s must be a valid UUID", fieldName)
	}

	return uid, nil
}

func ValidationRequired(fieldName string, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	return nil
}

func ValidationLength(fieldName, value string, min, max int) error {
	l := len(value)

	if l < min || l > max {
		return fmt.Errorf("%s must be between %d and %d characters", fieldName, min, max)
	}
	return nil
}

func ValidationRegex(fieldName, value string, reg *regexp.Regexp, msg string) error {
	if !reg.MatchString(value) {
		return fmt.Errorf("%s %s", fieldName, msg)
	}

	return nil
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
