package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go_pco_arrivals/internal/config"
	"go_pco_arrivals/internal/models"
	"go_pco_arrivals/internal/services"
	"go_pco_arrivals/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	config *config.Config
	db     *gorm.DB
	logger *utils.Logger
	auth   *services.AuthService
	pco    *services.PCOService
}

type LoginRequest struct {
	RememberMe bool `json:"remember_me"`
}

type AuthStatusResponse struct {
	IsAuthenticated bool                  `json:"is_authenticated"`
	User            *models.User          `json:"user,omitempty"`
	Session         *services.SessionData `json:"session,omitempty"`
	ExpiresAt       *time.Time            `json:"expires_at,omitempty"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewAuthHandler(config *config.Config, db *gorm.DB, logger *utils.Logger, auth *services.AuthService, pco *services.PCOService) *AuthHandler {
	return &AuthHandler{
		config: config,
		db:     db,
		logger: logger,
		auth:   auth,
		pco:    pco,
	}
}

// InitiateOAuth starts the OAuth flow by redirecting to PCO
func (h *AuthHandler) InitiateOAuth(c *fiber.Ctx) error {
	// Generate state and nonce for security
	state, err := h.auth.GenerateState()
	if err != nil {
		h.logger.Error("Failed to generate state", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initiate authentication",
		})
	}

	nonce, err := h.auth.GenerateNonce()
	if err != nil {
		h.logger.Error("Failed to generate nonce", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initiate authentication",
		})
	}

	// Store state and nonce in session for validation (TODO: implement proper session handling)
	// For now, we'll use cookies as a temporary solution
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   h.config.Server.Port == 443,
		SameSite: "Lax",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_remember_me",
		Value:    fmt.Sprintf("%t", c.Query("remember_me") == "true"),
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   h.config.Server.Port == 443,
		SameSite: "Lax",
	})

	// Generate authorization URL
	authURL := h.pco.GetAuthorizationURL(state, nonce)

	h.logger.Info("Initiating OAuth flow", "state", state, "nonce", nonce)
	return c.Redirect(authURL, http.StatusTemporaryRedirect)
}

// OAuthCallback handles the OAuth callback from PCO
func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
	// Get parameters from callback
	code := c.Query("code")
	state := c.Query("state")
	error := c.Query("error")
	errorDescription := c.Query("error_description")

	// Check for OAuth errors
	if error != "" {
		h.logger.Error("OAuth error", "error", error, "description", errorDescription)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":             "OAuth authentication failed",
			"error_description": errorDescription,
		})
	}

	// Validate required parameters
	if code == "" || state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required OAuth parameters",
		})
	}

	// Get session data from cookies (TODO: implement proper session handling)
	storedState := c.Cookies("oauth_state")
	rememberMeStr := c.Cookies("oauth_remember_me")
	rememberMe := rememberMeStr == "true"

	// Validate state parameter
	if storedState == "" || storedState != state {
		h.logger.Error("Invalid OAuth state", "received", state, "stored", storedState)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid OAuth state parameter",
		})
	}

	// Exchange code for access token
	authResp, err := h.pco.ExchangeCodeForToken(code)
	if err != nil {
		h.logger.Error("Failed to exchange code for token", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to complete authentication",
		})
	}

	// Get current user from PCO
	pcoUser, err := h.pco.GetCurrentUser(authResp.AccessToken)
	if err != nil {
		h.logger.Error("Failed to get current user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	// Validate user is authorized
	if !h.pco.ValidateUser(pcoUser.ID) {
		h.logger.Error("Unauthorized user attempted login", "pco_user_id", pcoUser.ID)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User is not authorized to access this application",
		})
	}

	// Create or update user in database
	user, err := h.pco.CreateOrUpdateUser(pcoUser, authResp.AccessToken)
	if err != nil {
		h.logger.Error("Failed to create/update user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save user information",
		})
	}

	// Update user tokens
	user.AccessToken = authResp.AccessToken
	user.RefreshToken = authResp.RefreshToken
	user.TokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second)
	user.UpdatedAt = time.Now()

	if err := h.db.Save(user).Error; err != nil {
		h.logger.Error("Failed to save user tokens", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save authentication tokens",
		})
	}

	// Create session
	sessionData, err := h.auth.CreateSession(user, rememberMe)
	if err != nil {
		h.logger.Error("Failed to create session", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user session",
		})
	}

	// Set session cookie with proper domain settings for cross-origin
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionData.Token,
		Expires:  sessionData.ExpiresAt,
		HTTPOnly: true,
		Secure:   false, // Set to false for localhost development
		SameSite: "Lax",
		Path:     "/",
	})

	// Clear OAuth session data
	c.ClearCookie("oauth_state")
	c.ClearCookie("oauth_nonce")
	c.ClearCookie("oauth_remember_me")

	h.logger.Info("User authenticated successfully", "user_id", user.ID, "pco_user_id", pcoUser.ID)

	// Redirect to dashboard or return success response
	if c.Get("Accept") == "application/json" {
		return c.JSON(fiber.Map{
			"success": true,
			"user":    user,
			"session": fiber.Map{
				"token":      sessionData.Token,
				"expires_at": sessionData.ExpiresAt,
			},
		})
	}

	// Redirect to frontend dashboard
	frontendURL := "http://localhost:5173/dashboard"
	return c.Redirect(frontendURL, http.StatusTemporaryRedirect)
}

// GetAuthStatus returns the current authentication status
func (h *AuthHandler) GetAuthStatus(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.JSON(AuthStatusResponse{
			IsAuthenticated: false,
		})
	}

	// Validate session
	sessionData, err := h.auth.ValidateSession(sessionToken)
	if err != nil {
		// Clear invalid session cookie
		c.ClearCookie("session_token")
		return c.JSON(AuthStatusResponse{
			IsAuthenticated: false,
		})
	}

	// Get user information
	user, err := h.auth.GetUserByID(sessionData.UserID)
	if err != nil {
		h.logger.Error("Failed to get user for auth status", "error", err, "user_id", sessionData.UserID)
		return c.JSON(AuthStatusResponse{
			IsAuthenticated: false,
		})
	}

	// Validate user access
	if err := h.auth.ValidateUserAccess(user); err != nil {
		h.logger.Error("User access validation failed", "error", err, "user_id", user.ID)
		return c.JSON(AuthStatusResponse{
			IsAuthenticated: false,
		})
	}

	return c.JSON(AuthStatusResponse{
		IsAuthenticated: true,
		User:            user,
		Session:         sessionData,
		ExpiresAt:       &sessionData.ExpiresAt,
	})
}

// Logout logs out the current user
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken != "" {
		// Revoke session
		if err := h.auth.RevokeSession(sessionToken); err != nil {
			h.logger.Error("Failed to revoke session", "error", err)
		}
	}

	// Clear session cookie
	c.ClearCookie("session_token")

	return c.JSON(LogoutResponse{
		Success: true,
		Message: "Successfully logged out",
	})
}

// RefreshToken refreshes the user's access token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No active session",
		})
	}

	// Validate session
	sessionData, err := h.auth.ValidateSession(sessionToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid session",
		})
	}

	// Get user
	user, err := h.auth.GetUserByID(sessionData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	// Refresh tokens
	if err := h.auth.RefreshUserTokens(user); err != nil {
		h.logger.Error("Failed to refresh user tokens", "error", err, "user_id", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to refresh authentication tokens",
		})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"message":    "Tokens refreshed successfully",
		"expires_at": user.TokenExpiry,
	})
}

// GetUserProfile returns the current user's profile
func (h *AuthHandler) GetUserProfile(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No active session",
		})
	}

	// Validate session
	sessionData, err := h.auth.ValidateSession(sessionToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid session",
		})
	}

	// Get user
	user, err := h.auth.GetUserByID(sessionData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	// Return user profile (excluding sensitive data)
	profile := fiber.Map{
		"id":            user.ID,
		"pco_user_id":   user.PCOUserID,
		"name":          user.Name,
		"email":         user.Email,
		"avatar":        user.Avatar,
		"is_admin":      user.IsAdmin,
		"is_active":     user.IsActive,
		"last_login":    user.LastLogin,
		"last_activity": user.LastActivity,
		"created_at":    user.CreatedAt,
	}

	return c.JSON(profile)
}

// UpdateUserProfile updates the current user's profile
func (h *AuthHandler) UpdateUserProfile(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No active session",
		})
	}

	// Validate session
	sessionData, err := h.auth.ValidateSession(sessionToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid session",
		})
	}

	// Get user
	user, err := h.auth.GetUserByID(sessionData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	// Parse update request
	var updateRequest struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}

	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update user fields
	if updateRequest.Name != "" {
		user.Name = updateRequest.Name
	}
	if updateRequest.Email != "" {
		user.Email = updateRequest.Email
	}
	if updateRequest.Avatar != "" {
		user.Avatar = updateRequest.Avatar
	}

	user.UpdatedAt = time.Now()

	// Save changes
	if err := h.db.Save(user).Error; err != nil {
		h.logger.Error("Failed to update user profile", "error", err, "user_id", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile updated successfully",
		"user": fiber.Map{
			"id":            user.ID,
			"pco_user_id":   user.PCOUserID,
			"name":          user.Name,
			"email":         user.Email,
			"avatar":        user.Avatar,
			"is_admin":      user.IsAdmin,
			"is_active":     user.IsActive,
			"last_login":    user.LastLogin,
			"last_activity": user.LastActivity,
			"created_at":    user.CreatedAt,
		},
	})
}
