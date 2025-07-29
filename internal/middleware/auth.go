package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AuthServiceInterface defines the interface for auth service methods used in middleware
type AuthServiceInterface interface {
	ValidateSessionForMiddleware(token string) (interface{}, error)
}

// RequireAuth middleware that validates session tokens and sets user_id in context
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the auth service from the app
		authService := c.Locals("auth_service")
		if authService == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Auth service not available",
			})
		}

		// Get session token from cookie
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No session token provided",
			})
		}

		// Validate session using auth service
		auth, ok := authService.(AuthServiceInterface)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid auth service type",
			})
		}

		sessionData, err := auth.ValidateSessionForMiddleware(sessionToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid session",
			})
		}

		// Extract user_id from SessionData struct
		if sessionDataStruct, ok := sessionData.(interface{ GetUserID() uint }); ok {
			c.Locals("user_id", sessionDataStruct.GetUserID())
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid session data",
		})
	}
}

// RequireAdmin middleware that requires admin privileges
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First check authentication
		if err := RequireAuth()(c); err != nil {
			return err
		}

		// Get the auth service from the app
		authService := c.Locals("auth_service")
		if authService == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Auth service not available",
			})
		}

		// Get session token from cookie
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No session token provided",
			})
		}

		// Validate session and check admin status
		auth, ok := authService.(AuthServiceInterface)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid auth service type",
			})
		}

		sessionData, err := auth.ValidateSessionForMiddleware(sessionToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid session",
			})
		}

		// Check if user is admin
		var isAdmin bool
		if sessionDataStruct, ok := sessionData.(interface{ IsAdmin() bool }); ok {
			isAdmin = sessionDataStruct.IsAdmin()
		} else if sessionDataStruct, ok := sessionData.(interface{ GetIsAdmin() bool }); ok {
			isAdmin = sessionDataStruct.GetIsAdmin()
		}

		if !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		return c.Next()
	}
}

// OptionalAuth middleware that optionally sets user_id if authenticated
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get session token from cookie
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.Next() // Continue without authentication
		}

		// Get the auth service from the app
		authService := c.Locals("auth_service")
		if authService == nil {
			return c.Next() // Continue without authentication
		}

		// Validate session using auth service
		auth, ok := authService.(AuthServiceInterface)
		if !ok {
			return c.Next() // Continue without authentication
		}

		sessionData, err := auth.ValidateSessionForMiddleware(sessionToken)
		if err != nil {
			return c.Next() // Continue without authentication
		}

		// Extract user_id from session data and set in context
		if sessionDataStruct, ok := sessionData.(interface{ GetUserID() uint }); ok {
			c.Locals("user_id", sessionDataStruct.GetUserID())
		}

		return c.Next()
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
