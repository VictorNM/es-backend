package user

import (
	"github.com/go-playground/validator/v10"
	"log"
	"unicode"
)

var validate = validator.New()

// register all custom validation logic here
func init() {
	err := validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return validatePassword(fl.Field().String())
	})

	if err != nil {
		log.Fatal(err)
	}
}

func validatePassword(password string) bool {
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	hasLetter := false
	hasDigit := false
	for _, r := range password {
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if unicode.IsLetter(r) {
			hasLetter = true
		}
	}

	return hasLetter && hasDigit
}
