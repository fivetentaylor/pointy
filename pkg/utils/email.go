package utils

import (
	"net/mail"
	"strings"
)

func NormalizeEmail(email string) string {
	email = strings.ToLower(email)             // Convert to lowercase
	email = strings.ReplaceAll(email, " ", "") // Remove all spaces
	return email
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
