package internal

import (
	"github.com/go-playground/validator/v10"
	"log"
	"unicode"
)

// Validate is the central place to Validate all input, using lib github.com/go-playground/validator tag
// in unit test, this function can be disable by replace with Validate = func(o interface{}) error {return nil}
var Validate = func(o interface{}) error {
	return v.Struct(o)
}

var v = validator.New()

// register all custom validation logic here
func init() {
	err := v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return validatePassword(fl.Field().String())
	})

	if err != nil {
		log.Fatal(err)
	}
}

func validatePassword(password string) bool {
	if len(password) < 8 {
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