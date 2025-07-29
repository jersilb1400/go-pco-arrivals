package types

import "github.com/gofiber/websocket/v2"

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

type WebSocketClient struct {
	Conn       *websocket.Conn
	ID         string
	LocationID string
	IsAdmin    bool
	UserID     string
}
