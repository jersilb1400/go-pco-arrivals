package services

import (
	"go_pco_arrivals/internal/utils"
)

type CleanupService struct {
	notificationService *NotificationService
	logger              *utils.Logger
	running             bool
}

func NewCleanupService(db interface{}, notificationService *NotificationService) *CleanupService {
	return &CleanupService{
		notificationService: notificationService,
		logger:              utils.NewLogger().WithComponent("cleanup_service"),
		running:             false,
	}
}

func (s *CleanupService) Start() {
	s.running = true
	s.logger.Info("Cleanup service started")

	// Cleanup service implementation will be added later
}

func (s *CleanupService) Stop() {
	s.running = false
	s.logger.Info("Cleanup service stopped")
}
