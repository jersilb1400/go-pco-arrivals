package handlers

import (
	"encoding/json"
	"go_pco_arrivals/internal/services"
	"go_pco_arrivals/internal/types"
	"go_pco_arrivals/internal/utils"

	"github.com/gofiber/websocket/v2"
)

type WebSocketHandler struct {
	hub         *services.WebSocketHub
	authService *services.AuthService
	logger      *utils.Logger
}

func NewWebSocketHandler(hub *services.WebSocketHub, authService *services.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
		logger:      utils.NewLogger().WithComponent("websocket_handler"),
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) {
	client := &types.WebSocketClient{
		Conn: c,
		ID:   utils.GenerateID(),
	}

	h.logger.Info("WebSocket client connected", "client_id", client.ID)

	// Register client with hub
	h.hub.Register(client)

	// Handle incoming messages
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			h.logger.Error("Failed to read WebSocket message", "error", err, "client_id", client.ID)
			break
		}

		h.handleMessage(client, message)
	}

	// Cleanup when connection closes
	h.hub.Unregister(client)
	h.logger.Info("WebSocket client disconnected", "client_id", client.ID)
}

func (h *WebSocketHandler) HandleBillboardWebSocket(c *websocket.Conn) {
	// Extract location ID from URL path
	locationID := c.Params("locationId")
	if locationID == "" {
		locationID = "all"
	}

	client := &types.WebSocketClient{
		Conn:       c,
		ID:         utils.GenerateID(),
		LocationID: locationID,
	}

	h.logger.Info("Billboard WebSocket client connected",
		"client_id", client.ID,
		"location_id", locationID)

	// Register client with hub
	h.hub.Register(client)

	// Send initial connection confirmation
	h.sendMessage(client, types.WebSocketMessage{
		Type: "connection_established",
		Data: map[string]interface{}{
			"client_id":   client.ID,
			"location_id": locationID,
		},
		Timestamp: utils.GetCurrentTimestamp(),
	})

	// Handle incoming messages
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			h.logger.Error("Failed to read billboard WebSocket message",
				"error", err,
				"client_id", client.ID)
			break
		}

		h.handleMessage(client, message)
	}

	// Cleanup when connection closes
	h.hub.Unregister(client)
	h.logger.Info("Billboard WebSocket client disconnected",
		"client_id", client.ID,
		"location_id", locationID)
}

func (h *WebSocketHandler) handleMessage(client *types.WebSocketClient, message []byte) {
	var wsMessage types.WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		h.logger.Error("Failed to unmarshal WebSocket message", "error", err)
		return
	}

	h.logger.Debug("Received WebSocket message",
		"type", wsMessage.Type,
		"client_id", client.ID)

	switch wsMessage.Type {
	case "ping":
		// Respond to ping with pong
		h.sendMessage(client, types.WebSocketMessage{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": wsMessage.Data.(map[string]interface{})["timestamp"],
			},
			Timestamp: utils.GetCurrentTimestamp(),
		})

	case "subscribe_location":
		if data, ok := wsMessage.Data.(map[string]interface{}); ok {
			if locationID, ok := data["location_id"].(string); ok {
				client.LocationID = locationID
				h.logger.Info("Client subscribed to location",
					"client_id", client.ID,
					"location_id", locationID)
			}
		}

	case "subscribe_notifications":
		// Client wants to receive notification updates
		h.logger.Info("Client subscribed to notifications", "client_id", client.ID)

	case "subscribe_billboard_state":
		// Client wants to receive billboard state updates
		h.logger.Info("Client subscribed to billboard state", "client_id", client.ID)

	default:
		h.logger.Warn("Unknown WebSocket message type", "type", wsMessage.Type)
	}
}

func (h *WebSocketHandler) sendMessage(client *types.WebSocketClient, message types.WebSocketMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal WebSocket message", "error", err)
		return
	}

	if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		h.logger.Error("Failed to send WebSocket message", "error", err, "client_id", client.ID)
	}
}

func (h *WebSocketHandler) BroadcastToLocation(locationID string, messageType string, data interface{}) {
	message := types.WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: utils.GetCurrentTimestamp(),
	}

	h.hub.BroadcastToLocation(locationID, messageType, message)
}

func (h *WebSocketHandler) BroadcastToAdmins(messageType string, data interface{}) {
	message := types.WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: utils.GetCurrentTimestamp(),
	}

	h.hub.BroadcastToAdmins(messageType, message)
}

func (h *WebSocketHandler) Broadcast(messageType string, data interface{}) {
	message := types.WebSocketMessage{
		Type:      messageType,
		Data:      data,
		Timestamp: utils.GetCurrentTimestamp(),
	}

	h.hub.Broadcast(messageType, message)
}
