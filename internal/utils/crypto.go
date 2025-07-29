package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"time"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateSecureToken() (string, error) {
	return GenerateRandomString(32)
}

// GenerateID creates a random 16-character hex string
func GenerateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetCurrentTimestamp returns the current timestamp in ISO format
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
