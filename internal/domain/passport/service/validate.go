package service

import (
	"strings"
	"unicode"
	"unicode/utf8"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/Flyskea/gotools/errors"

	"loginhub/internal/base/reason"
)

const (
	maxNameLength     = 20
	minPasswordLength = 8
	maxPasswordLength = 20
)

var (
	specialChInPassword = ".@$!%*#_~?&^"
)

var (
	verifier = emailverifier.NewVerifier()
)

func validateUserName(name string) error {
	// check length
	if utf8.RuneCountInString(name) > maxNameLength {
		return errors.BadRequest(reason.UserInvalidName)
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return errors.BadRequest(reason.UserInvalidName)
		}
	}
	return nil
}

func validatePassword(password string) error {
	// check length
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return errors.BadRequest(reason.UserInvalidPassword)
	}

	hasLetter := false
	hasDigit := false
	hasSpecial := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if strings.ContainsRune(specialChInPassword, char) {
			hasSpecial = true
		}
	}
	if !hasLetter || !hasDigit || !hasSpecial {
		return errors.BadRequest(reason.UserInvalidPassword)
	}

	return nil
}

func validateEmail(email string) error {
	result, err := verifier.Verify(email)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	if !result.Syntax.Valid || !result.HasMxRecords {
		return errors.BadRequest(reason.InvalidEmail)
	}
	return nil
}
