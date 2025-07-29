package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func SanitizeString(input string) string {
	// Remove any potentially dangerous characters
	re := regexp.MustCompile(`[<>\"'&]`)
	return re.ReplaceAllString(input, "")
}

func ValidatePCOUserID(userID string) bool {
	// PCO user IDs are typically alphanumeric
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(userID)
}

func ValidateSecurityCode(code string) bool {
	// Security codes are typically 3-6 alphanumeric characters
	re := regexp.MustCompile(`^[A-Z0-9]{3,6}$`)
	return re.MatchString(strings.ToUpper(code))
}

func ValidateLocationID(locationID string) bool {
	// Location IDs are typically alphanumeric
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(locationID)
}
