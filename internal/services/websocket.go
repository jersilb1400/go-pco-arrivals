package services

import (
	"encoding/json"
	"go_pco_arrivals/internal/types"
	"go_pco_arrivals/internal/utils"
	"sync"
)

type WebSocketHub struct {
	logger    *utils.Logger
	running   bool
	clients   map[string]*types.WebSocketClient
	admins    map[string]*types.WebSocketClient
	locations map[string]map[string]*types.WebSocketClient
	mutex     sync.RWMutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		logger:    utils.NewLogger().WithComponent("websocket_hub"),
		running:   false,
		clients:   make(map[string]*types.WebSocketClient),
		admins:    make(map[string]*types.WebSocketClient),
		locations: make(map[string]map[string]*types.WebSocketClient),
	}
}

func (h *WebSocketHub) Run() {
	h.running = true
	h.logger.Info("WebSocket hub started")
}

func (h *WebSocketHub) Stop() {
	h.running = false
	h.logger.Info("WebSocket hub stopped")
}

func (h *WebSocketHub) Register(client *types.WebSocketClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client.ID] = client

	// Register for location-specific updates
	if client.LocationID != "" {
		if h.locations[client.LocationID] == nil {
			h.locations[client.LocationID] = make(map[string]*types.WebSocketClient)
		}
		h.locations[client.LocationID][client.ID] = client
	}

	// Register admin clients
	if client.IsAdmin {
		h.admins[client.ID] = client
	}

	h.logger.Info("Client registered",
		"client_id", client.ID,
		"location_id", client.LocationID,
		"is_admin", client.IsAdmin)
}

func (h *WebSocketHub) Unregister(client *types.WebSocketClient) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove from general clients
	delete(h.clients, client.ID)

	// Remove from admin clients
	delete(h.admins, client.ID)

	// Remove from location-specific clients
	if client.LocationID != "" {
		if locationClients, exists := h.locations[client.LocationID]; exists {
			delete(locationClients, client.ID)
			if len(locationClients) == 0 {
				delete(h.locations, client.LocationID)
			}
		}
	}

	h.logger.Info("Client unregistered",
		"client_id", client.ID,
		"location_id", client.LocationID)
}

func (h *WebSocketHub) Broadcast(messageType string, data interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	message := types.WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: utils.GetCurrentTimestamp(),
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal broadcast message", "error", err)
		return
	}

	for _, client := range h.clients {
		if err := client.Conn.WriteMessage(1, messageData); err != nil {
			h.logger.Error("Failed to send broadcast message to client",
				"error", err,
				"client_id", client.ID)
		}
	}

	h.logger.Debug("Broadcast message sent",
		"type", messageType,
		"recipients", len(h.clients))
}

func (h *WebSocketHub) BroadcastToLocation(locationID string, messageType string, data interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	message := types.WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: utils.GetCurrentTimestamp(),
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal location broadcast message", "error", err)
		return
	}

	locationClients, exists := h.locations[locationID]
	if !exists {
		h.logger.Debug("No clients for location", "location_id", locationID)
		return
	}

	for _, client := range locationClients {
		if err := client.Conn.WriteMessage(1, messageData); err != nil {
			h.logger.Error("Failed to send location broadcast message to client",
				"error", err,
				"client_id", client.ID,
				"location_id", locationID)
		}
	}

	h.logger.Debug("Location broadcast message sent",
		"type", messageType,
		"location_id", locationID,
		"recipients", len(locationClients))
}

func (h *WebSocketHub) BroadcastToAdmins(messageType string, data interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	message := types.WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: utils.GetCurrentTimestamp(),
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal admin broadcast message", "error", err)
		return
	}

	for _, client := range h.admins {
		if err := client.Conn.WriteMessage(1, messageData); err != nil {
			h.logger.Error("Failed to send admin broadcast message to client",
				"error", err,
				"client_id", client.ID)
		}
	}

	h.logger.Debug("Admin broadcast message sent",
		"type", messageType,
		"recipients", len(h.admins))
}

func (h *WebSocketHub) GetStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	locationStats := make(map[string]int)
	for locationID, clients := range h.locations {
		locationStats[locationID] = len(clients)
	}

	return map[string]interface{}{
		"running":         h.running,
		"total_clients":   len(h.clients),
		"admin_clients":   len(h.admins),
		"location_stats":  locationStats,
		"total_locations": len(h.locations),
	}
}
