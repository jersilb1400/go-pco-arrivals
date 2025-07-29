package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"go_pco_arrivals/internal/config"
	"go_pco_arrivals/internal/models"
	"go_pco_arrivals/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	config *config.Config
	db     *gorm.DB
	logger *utils.Logger
	pco    *PCOService
}

type Claims struct {
	UserID    uint   `json:"user_id"`
	PCOUserID string `json:"pco_user_id"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type SessionData struct {
	UserID       uint      `json:"user_id"`
	PCOUserID    string    `json:"pco_user_id"`
	Email        string    `json:"email"`
	IsAdmin      bool      `json:"is_admin"`
	IsRememberMe bool      `json:"is_remember_me"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GetUserID returns the user ID
func (s *SessionData) GetUserID() uint {
	return s.UserID
}

// GetIsAdmin returns whether the user is an admin
func (s *SessionData) GetIsAdmin() bool {
	return s.IsAdmin
}

func NewAuthService(config *config.Config, db *gorm.DB, logger *utils.Logger, pco *PCOService) *AuthService {
	return &AuthService{
		config: config,
		db:     db,
		logger: logger,
		pco:    pco,
	}
}

// GenerateState generates a random state parameter for OAuth
func (s *AuthService) GenerateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateNonce generates a random nonce parameter for OAuth
func (s *AuthService) GenerateNonce() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// CreateSession creates a new session for a user
func (s *AuthService) CreateSession(user *models.User, isRememberMe bool) (*models.Session, error) {
	sessionToken, err := utils.GenerateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	var expiresAt time.Time
	if isRememberMe {
		expiresAt = time.Now().AddDate(0, 0, s.config.Auth.RememberMeDays)
	} else {
		expiresAt = time.Now().Add(time.Duration(s.config.Auth.SessionTTL) * time.Second)
	}

	session := &models.Session{
		UserID:       user.ID,
		Token:        sessionToken,
		IsRememberMe: isRememberMe,
		ExpiresAt:    expiresAt,
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// ValidateSession validates a session token and returns session data
func (s *AuthService) ValidateSession(token string) (*SessionData, error) {
	var session models.Session
	result := s.db.Preload("User").Where("token = ? AND expires_at > ?", token, time.Now()).First(&session)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired session")
		}
		return nil, fmt.Errorf("failed to validate session: %w", result.Error)
	}

	// Update last activity
	session.LastActivity = time.Now()
	s.db.Save(&session)

	return &SessionData{
		UserID:       session.User.ID,
		PCOUserID:    session.User.PCOUserID,
		Email:        session.User.Email,
		IsAdmin:      session.User.IsAdmin,
		IsRememberMe: session.IsRememberMe,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}

// CreateRememberMeSession creates a long-term session for "Remember Me" functionality
func (s *AuthService) CreateRememberMeSession(user *models.User) (*models.Session, error) {
	return s.CreateSession(user, true)
}

// GenerateJWT generates a JWT token for a user
func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:    user.ID,
		PCOUserID: user.PCOUserID,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go_pco_arrivals",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Auth.JWTSecret))
}

// ValidateJWT validates a JWT token and returns claims
func (s *AuthService) ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}

// RevokeSession revokes a session by deleting it from the database
func (s *AuthService) RevokeSession(token string) error {
	result := s.db.Where("token = ?", token).Delete(&models.Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to revoke session: %w", result.Error)
	}
	return nil
}

// RevokeAllUserSessions revokes all sessions for a specific user
func (s *AuthService) RevokeAllUserSessions(userID uint) error {
	result := s.db.Where("user_id = ?", userID).Delete(&models.Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to revoke user sessions: %w", result.Error)
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *AuthService) CleanupExpiredSessions() error {
	result := s.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", result.Error)
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return &user, nil
}

// GetUserByPCOID retrieves a user by PCO user ID
func (s *AuthService) GetUserByPCOID(pcoUserID string) (*models.User, error) {
	var user models.User
	result := s.db.Where("pco_user_id = ?", pcoUserID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return &user, nil
}

// UpdateUserLastActivity updates the last activity timestamp for a user
func (s *AuthService) UpdateUserLastActivity(userID uint) error {
	result := s.db.Model(&models.User{}).Where("id = ?", userID).Update("last_activity", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to update user activity: %w", result.Error)
	}
	return nil
}

// IsTokenExpiringSoon checks if a user's access token is expiring soon
func (s *AuthService) IsTokenExpiringSoon(user *models.User) bool {
	if user.TokenExpiry.IsZero() {
		return false
	}

	threshold := time.Duration(s.config.Auth.TokenRefreshThreshold) * time.Second
	return time.Until(user.TokenExpiry) < threshold
}

// RefreshUserTokens refreshes a user's PCO access tokens
func (s *AuthService) RefreshUserTokens(user *models.User) error {
	if user.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	authResp, err := s.pco.RefreshAccessToken(user.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to refresh access token: %w", err)
	}

	// Update user tokens
	user.AccessToken = authResp.AccessToken
	user.RefreshToken = authResp.RefreshToken
	user.TokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second)
	user.UpdatedAt = time.Now()

	if err := s.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to save refreshed tokens: %w", err)
	}

	return nil
}

// ValidateUserAccess validates that a user has valid access to the system
func (s *AuthService) ValidateUserAccess(user *models.User) error {
	// Check if user is active
	if !user.IsActive {
		return fmt.Errorf("user account is inactive")
	}

	// Check if user is authorized
	if !s.pco.ValidateUser(user.PCOUserID) {
		return fmt.Errorf("user is not authorized")
	}

	// Check if token needs refresh
	if s.IsTokenExpiringSoon(user) {
		if err := s.RefreshUserTokens(user); err != nil {
			return fmt.Errorf("failed to refresh tokens: %w", err)
		}
	}

	return nil
}

// ValidateSessionForMiddleware implements the middleware interface
// This wraps the existing ValidateSession method to return interface{} instead of *SessionData
func (s *AuthService) ValidateSessionForMiddleware(token string) (interface{}, error) {
	sessionData, err := s.ValidateSession(token)
	if err != nil {
		return nil, err
	}
	return sessionData, nil
}
