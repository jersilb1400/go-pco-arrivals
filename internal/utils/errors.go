package utils

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidInput       = errors.New("invalid input")
	ErrDatabaseError      = errors.New("database error")
	ErrPCOAPIError        = errors.New("PCO API error")
	ErrWebSocketError     = errors.New("WebSocket error")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
)
