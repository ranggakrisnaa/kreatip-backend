package websocket

import (
	"context"
	"sync"

	"github.com/gofiber/fiber/v2"
	fws "github.com/gofiber/websocket/v2"
	"github.com/kreatip/kreatip-backend/internal/usecase"
	"github.com/sirupsen/logrus"
)

// Hub manages active WebSocket connections keyed by userID.
// In-memory for MVP; replace with Redis Pub/Sub for multi-instance.
type Hub struct {
	mu      sync.RWMutex
	clients map[string][]*fws.Conn // userID → connections
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string][]*fws.Conn)}
}

func (h *Hub) Register(userID string, conn *fws.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], conn)
}

func (h *Hub) Unregister(userID string, conn *fws.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[userID]
	for i, c := range conns {
		if c == conn {
			h.clients[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func (h *Hub) Broadcast(userID string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, conn := range h.clients[userID] {
		_ = conn.WriteMessage(1, msg) // 1 = TextMessage
	}
}

// AlertHandler handles WS connections for OBS overlay.
type AlertHandler struct {
	Hub       *Hub
	AlertUC   usecase.AlertUseCase
	Log       *logrus.Logger
}

func NewAlertHandler(hub *Hub, uc usecase.AlertUseCase, log *logrus.Logger) *AlertHandler {
	return &AlertHandler{Hub: hub, AlertUC: uc, Log: log}
}

// Upgrade upgrades the HTTP connection to WebSocket.
// Route: GET /ws/alert?token=xxx
func (h *AlertHandler) Upgrade(c *fiber.Ctx) error {
	if !fws.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}
	return c.Next()
}

// Handle is the actual WS handler.
func (h *AlertHandler) Handle(c *fws.Conn) {
	token := c.Query("token")
	if token == "" {
		_ = c.Close()
		return
	}

	userID, err := h.AlertUC.AuthorizeToken(context.Background(), token)
	if err != nil || userID == "" {
		_ = c.Close()
		return
	}

	h.Hub.Register(userID, c)
	defer h.Hub.Unregister(userID, c)

	// Keep alive — read until disconnect
	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}
