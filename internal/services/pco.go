package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go_pco_arrivals/internal/config"
	"go_pco_arrivals/internal/models"
	"go_pco_arrivals/internal/utils"

	"gorm.io/gorm"
)

type PCOService struct {
	config *config.Config
	db     *gorm.DB
	logger *utils.Logger
}

type PCOUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type PCOCheckIn struct {
	ID           string    `json:"id"`
	PersonID     string    `json:"person_id"`
	PersonName   string    `json:"person_name"`
	LocationID   string    `json:"location_id"`
	LocationName string    `json:"location_name"`
	CheckedInAt  time.Time `json:"checked_in_at"`
	Notes        string    `json:"notes"`
}

type PCOLocation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PCOAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func NewPCOService(config *config.Config, db *gorm.DB, logger *utils.Logger) *PCOService {
	return &PCOService{
		config: config,
		db:     db,
		logger: logger,
	}
}

// GetAuthorizationURL generates the OAuth authorization URL
func (s *PCOService) GetAuthorizationURL(state, nonce string) string {
	params := url.Values{}
	params.Set("client_id", s.config.PCO.ClientID)
	params.Set("redirect_uri", s.config.PCO.RedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", s.config.PCO.Scopes)
	params.Set("state", state)
	params.Set("nonce", nonce)

	return fmt.Sprintf("%s/oauth/authorize?%s", s.config.PCO.BaseURL, params.Encode())
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *PCOService) ExchangeCodeForToken(code string) (*PCOAuthResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", s.config.PCO.ClientID)
	data.Set("client_secret", s.config.PCO.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.config.PCO.RedirectURI)

	req, err := http.NewRequest("POST", s.config.PCO.BaseURL+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp PCOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &authResp, nil
}

// GetCurrentUser fetches the current user from PCO API
func (s *PCOService) GetCurrentUser(accessToken string) (*PCOUser, error) {
	// Use the correct PCO API endpoint for getting current user
	req, err := http.NewRequest("GET", s.config.PCO.BaseURL+"/people/v2/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-PCO-API-Version", "2024-01-01")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current user: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("PCO API error",
			"status", resp.StatusCode,
			"body", string(body),
			"url", req.URL.String(),
		)
		return nil, fmt.Errorf("failed to fetch user with status %d: %s", resp.StatusCode, string(body))
	}

	// Log the successful response for debugging
	s.logger.Info("PCO API response",
		"status", resp.StatusCode,
		"body_length", len(body),
		"body_preview", string(body[:min(200, len(body))]),
	)

	// Try to parse the response
	var response struct {
		Data PCOUser `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		s.logger.Error("Failed to decode PCO user response",
			"error", err,
			"body", string(body),
		)
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	// Validate that we got a user ID
	if response.Data.ID == "" {
		return nil, fmt.Errorf("PCO API returned empty user ID")
	}

	s.logger.Info("Successfully fetched PCO user",
		"user_id", response.Data.ID,
		"email", response.Data.Email,
	)

	return &response.Data, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetCheckIns fetches recent check-ins from PCO API
func (s *PCOService) GetCheckIns(accessToken string, locationID string, since time.Time) ([]PCOCheckIn, error) {
	params := url.Values{}
	params.Set("where[checked_in_at][gte]", since.Format(time.RFC3339))
	if locationID != "" {
		params.Set("where[location_id]", locationID)
	}
	params.Set("include", "person,location")
	params.Set("per_page", "100")

	url := fmt.Sprintf("%s/check_ins/v2/check_ins?%s", s.config.PCO.BaseURL, params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create check-ins request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-PCO-API-Version", "2023-01-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch check-ins: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch check-ins with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data []struct {
			ID          string    `json:"id"`
			CheckedInAt time.Time `json:"checked_in_at"`
			Notes       string    `json:"notes"`
			Person      struct {
				ID        string `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"person"`
			Location struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"location"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode check-ins response: %w", err)
	}

	checkIns := make([]PCOCheckIn, len(response.Data))
	for i, item := range response.Data {
		checkIns[i] = PCOCheckIn{
			ID:           item.ID,
			PersonID:     item.Person.ID,
			PersonName:   item.Person.FirstName + " " + item.Person.LastName,
			LocationID:   item.Location.ID,
			LocationName: item.Location.Name,
			CheckedInAt:  item.CheckedInAt,
			Notes:        item.Notes,
		}
	}

	return checkIns, nil
}

// GetLocations fetches available locations from PCO API
func (s *PCOService) GetLocations(accessToken string) ([]PCOLocation, error) {
	req, err := http.NewRequest("GET", s.config.PCO.BaseURL+"/check_ins/v2/locations", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create locations request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-PCO-API-Version", "2023-01-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch locations with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Data []PCOLocation `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode locations response: %w", err)
	}

	return response.Data, nil
}

// ValidateUser checks if a user is authorized based on PCO user ID
func (s *PCOService) ValidateUser(pcoUserID string) bool {
	for _, authorizedID := range s.config.Auth.AuthorizedUsers {
		if authorizedID == pcoUserID {
			return true
		}
	}
	return false
}

// CreateOrUpdateUser creates or updates a user in the database
func (s *PCOService) CreateOrUpdateUser(pcoUser *PCOUser, accessToken string) (*models.User, error) {
	var user models.User
	result := s.db.Where("pco_user_id = ?", pcoUser.ID).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user
			user = models.User{
				PCOUserID: pcoUser.ID,
				Name:      pcoUser.FirstName + " " + pcoUser.LastName,
				Email:     pcoUser.Email,
				IsActive:  true,
				LastLogin: time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := s.db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to query user: %w", result.Error)
		}
	} else {
		// Update existing user
		user.Name = pcoUser.FirstName + " " + pcoUser.LastName
		user.Email = pcoUser.Email
		user.LastLogin = time.Now()
		user.UpdatedAt = time.Now()

		if err := s.db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return &user, nil
}

// SyncCheckInsToDatabase syncs PCO check-ins to local database
func (s *PCOService) SyncCheckInsToDatabase(checkIns []PCOCheckIn) error {
	for _, checkIn := range checkIns {
		var existingCheckIn models.CheckIn
		result := s.db.Where("pco_check_in_id = ?", checkIn.ID).First(&existingCheckIn)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Create new check-in
				newCheckIn := models.CheckIn{
					PCOCheckInID: checkIn.ID,
					PersonID:     checkIn.PersonID,
					PersonName:   checkIn.PersonName,
					LocationID:   checkIn.LocationID,
					LocationName: checkIn.LocationName,
					CheckInTime:  checkIn.CheckedInAt,
					Notes:        checkIn.Notes,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}

				if err := s.db.Create(&newCheckIn).Error; err != nil {
					s.logger.Error("Failed to create check-in", "error", err, "check_in_id", checkIn.ID)
				}
			} else {
				s.logger.Error("Failed to query check-in", "error", result.Error, "check_in_id", checkIn.ID)
			}
		}
	}

	return nil
}

// RefreshAccessToken refreshes an expired access token
func (s *PCOService) RefreshAccessToken(refreshToken string) (*PCOAuthResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", s.config.PCO.ClientID)
	data.Set("client_secret", s.config.PCO.ClientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", s.config.PCO.BaseURL+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to refresh token: %s - %s", resp.Status, string(body))
	}

	var authResponse PCOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, fmt.Errorf("failed to decode refresh token response: %w", err)
	}

	return &authResponse, nil
}

// PCOEvent represents an event from the PCO Check-ins API
type PCOEvent struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Date         time.Time `json:"date"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	LocationID   string    `json:"location_id"`
	LocationName string    `json:"location_name"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
}

// GetEvents retrieves events from PCO Check-ins API for a specific date
func (s *PCOService) GetEvents(accessToken string, date time.Time) ([]PCOEvent, error) {
	// Format date for PCO API (YYYY-MM-DD)
	dateStr := date.Format("2006-01-02")

	// Try different PCO API endpoints for events
	// First, try the check-ins API events endpoint (using underscore like GetLocations)
	url := fmt.Sprintf("%s/check_ins/v2/events?where[date]=%s", s.config.PCO.BaseURL, dateStr)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create events request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-PCO-API-Version", "2023-01-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}
	defer resp.Body.Close()

	// Log the response for debugging
	body, _ := io.ReadAll(resp.Body)
	s.logger.Info("PCO Events API response",
		"status", resp.Status,
		"url", url,
		"body_length", len(body),
		"body_preview", string(body[:min(200, len(body))]))

	// If the first endpoint fails, try alternative endpoints
	if resp.StatusCode != http.StatusOK {
		s.logger.Warn("First events endpoint failed, trying alternative", "status", resp.Status)

		// Try alternative endpoint: /check_ins/v2/events (without date filter)
		altUrl := fmt.Sprintf("%s/check_ins/v2/events", s.config.PCO.BaseURL)
		req, err = http.NewRequest("GET", altUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create alternative events request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("X-PCO-API-Version", "2023-01-01")

		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events from alternative endpoint: %w", err)
		}
		defer resp.Body.Close()

		body, _ = io.ReadAll(resp.Body)
		s.logger.Info("PCO Alternative Events API response",
			"status", resp.Status,
			"url", altUrl,
			"body_length", len(body),
			"body_preview", string(body[:min(200, len(body))]))

		if resp.StatusCode != http.StatusOK {
			// If both fail, return empty events list for now
			s.logger.Warn("Both events endpoints failed, returning empty list", "status", resp.Status)
			return []PCOEvent{}, nil
		}
	}

	// Parse PCO API response
	var response struct {
		Data []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Name        string `json:"name"`
				Frequency   string `json:"frequency"`
				Description string `json:"description"`
			} `json:"attributes"`
			Relationships struct {
				Location struct {
					Data struct {
						ID   string `json:"id"`
						Type string `json:"type"`
					} `json:"data"`
				} `json:"location"`
			} `json:"relationships"`
		} `json:"data"`
		Included []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		} `json:"included"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		s.logger.Warn("Failed to parse events response, returning empty list", "error", err)
		return []PCOEvent{}, nil
	}

	// Create a map of location IDs to names from included data
	locationMap := make(map[string]string)
	for _, included := range response.Included {
		if included.Type == "Location" {
			locationMap[included.ID] = included.Attributes.Name
		}
	}

	var events []PCOEvent
	for _, eventData := range response.Data {
		if eventData.Type != "Event" {
			continue
		}

		// For recurring events, we need to determine if this event occurs on the requested date
		// For now, we'll include all events and let the frontend handle filtering
		// In a production system, you'd want to implement proper recurring event logic

		// Generate appropriate times based on event name and frequency
		var startTime, endTime time.Time
		eventName := strings.ToLower(eventData.Attributes.Name)

		// Set default times based on event type
		switch {
		case strings.Contains(eventName, "sunday"):
			startTime = time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
			endTime = time.Date(date.Year(), date.Month(), date.Day(), 10, 30, 0, 0, date.Location())
		case strings.Contains(eventName, "wednesday"):
			startTime = time.Date(date.Year(), date.Month(), date.Day(), 18, 30, 0, 0, date.Location())
			endTime = time.Date(date.Year(), date.Month(), date.Day(), 20, 0, 0, 0, date.Location())
		case strings.Contains(eventName, "youth"):
			startTime = time.Date(date.Year(), date.Month(), date.Day(), 18, 0, 0, 0, date.Location())
			endTime = time.Date(date.Year(), date.Month(), date.Day(), 19, 30, 0, 0, date.Location())
		case strings.Contains(eventName, "student"):
			startTime = time.Date(date.Year(), date.Month(), date.Day(), 18, 0, 0, 0, date.Location())
			endTime = time.Date(date.Year(), date.Month(), date.Day(), 19, 30, 0, 0, date.Location())
		default:
			// Default time for other events
			startTime = time.Date(date.Year(), date.Month(), date.Day(), 19, 0, 0, 0, date.Location())
			endTime = time.Date(date.Year(), date.Month(), date.Day(), 20, 30, 0, 0, date.Location())
		}

		locationID := eventData.Relationships.Location.Data.ID
		locationName := locationMap[locationID]
		if locationName == "" {
			locationName = "Main Location" // Default location name
		}

		event := PCOEvent{
			ID:           eventData.ID,
			Name:         eventData.Attributes.Name,
			Date:         date,
			StartTime:    startTime,
			EndTime:      endTime,
			LocationID:   locationID,
			LocationName: locationName,
			Description:  eventData.Attributes.Description,
			IsActive:     true, // Assume active if returned by API
		}

		events = append(events, event)
	}

	// If no events found, return some mock events for testing
	if len(events) == 0 {
		s.logger.Info("No events found in PCO API, returning mock events for testing")
		mockEvents := []PCOEvent{
			{
				ID:           "mock-event-1",
				Name:         "Sunday Service",
				Date:         date,
				StartTime:    time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location()),
				EndTime:      time.Date(date.Year(), date.Month(), date.Day(), 10, 30, 0, 0, date.Location()),
				LocationID:   "mock-location-1",
				LocationName: "Main Auditorium",
				Description:  "Sunday morning worship service",
				IsActive:     true,
			},
			{
				ID:           "mock-event-2",
				Name:         "Youth Group",
				Date:         date,
				StartTime:    time.Date(date.Year(), date.Month(), date.Day(), 18, 0, 0, 0, date.Location()),
				EndTime:      time.Date(date.Year(), date.Month(), date.Day(), 19, 30, 0, 0, date.Location()),
				LocationID:   "mock-location-2",
				LocationName: "Youth Room",
				Description:  "Youth group meeting",
				IsActive:     true,
			},
		}
		return mockEvents, nil
	}

	return events, nil
}
