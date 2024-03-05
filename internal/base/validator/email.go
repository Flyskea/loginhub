package validator

import (
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/go-playground/validator/v10"
)

var emailValidator = emailverifier.NewVerifier()

func ValidateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	result, err := emailValidator.Verify(email)
	if err != nil {
		return false
	}
	if !result.Syntax.Valid || !result.HasMxRecords {
		return false
	}
	return true
}
