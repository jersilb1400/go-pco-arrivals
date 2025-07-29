package services

import (
	"go_pco_arrivals/internal/utils"
)

type RealtimeService struct {
	logger *utils.Logger
}

func NewRealtimeService() *RealtimeService {
	return &RealtimeService{
		logger: utils.NewLogger().WithComponent("realtime_service"),
	}
}

func (s *RealtimeService) GetCheckIns(accessToken string, eventID string) ([]interface{}, error) {
	// Get check-ins will be implemented later
	return nil, nil
}

func (s *RealtimeService) GetCheckInsByLocation(accessToken string, locationID string) ([]interface{}, error) {
	// Get check-ins by location will be implemented later
	return nil, nil
}

func (s *RealtimeService) GetCheckInsByEvent(accessToken string, eventID string) ([]interface{}, error) {
	// Get check-ins by event will be implemented later
	return nil, nil
}

func (s *RealtimeService) GetLocations(accessToken string) ([]interface{}, error) {
	// Get locations will be implemented later
	return nil, nil
}

func (s *RealtimeService) GetLocation(accessToken string, locationID string) (interface{}, error) {
	// Get location will be implemented later
	return nil, nil
}
