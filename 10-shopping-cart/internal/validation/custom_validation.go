package validation

import (
	"log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/qwenhn/gin-restful-api/10-shopping-cart/internal/utils"
)

func RegisterCustomValidation(v *validator.Validate) {
	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	var blockedDomains = map[string]bool{
		"blacklist.com": true,
		"edu.vn":        true,
		"abc.com":       true,
	}
	v.RegisterValidation("email_advanced", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()

		parts := strings.Split(email, "@")
		log.Println("parts:", parts)
		if len(parts) != 2 {
			return false
		}

		domain := utils.NormalizeString(parts[1])

		return !blockedDomains[domain]
	})

	v.RegisterValidation("password_strong", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 8 {
			return false
		}

		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>?/\\|]`).MatchString(password)

		return hasLower && hasUpper && hasDigit && hasSpecial
	})
}
