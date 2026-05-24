package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Message struct {
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type Connection struct {
	ID     string
	Send   chan []byte
	logger *zap.Logger
	mu     sync.Mutex
	closed bool
}

func NewConnection(id string, logger *zap.Logger) *Connection {
	return &Connection{
		ID:     id,
		Send:   make(chan []byte, 256),
		logger: logger,
	}
}

func (c *Connection) SendJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	select {
	case c.Send <- data:
	default:
		c.logger.Warn("websocket send buffer full", zap.String("conn_id", c.ID))
	}

	return nil
}

func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}
	c.closed = true
	close(c.Send)
}

type Hub struct {
	connections map[string]*Connection
	mu          sync.RWMutex
	logger      *zap.Logger
}

func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		connections: make(map[string]*Connection),
		logger:      logger,
	}
}

func (h *Hub) Register(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[conn.ID] = conn
	h.logger.Info("websocket client connected", zap.String("conn_id", conn.ID))
}

func (h *Hub) Unregister(conn *Connection) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.connections[conn.ID]; ok {
		delete(h.connections, conn.ID)
		conn.Close()
		h.logger.Info("websocket client disconnected", zap.String("conn_id", conn.ID))
	}
}

func (h *Hub) Broadcast(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		h.logger.Error("broadcast marshal error", zap.Error(err))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conn := range h.connections {
		select {
		case conn.Send <- data:
		default:
			h.logger.Warn("broadcast buffer full", zap.String("conn_id", conn.ID))
		}
	}
}

func (h *Hub) SendTo(connID string, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		h.logger.Error("send marshal error", zap.Error(err))
		return
	}

	h.mu.RLock()
	conn, ok := h.connections[connID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	select {
	case conn.Send <- data:
	case <-time.After(time.Second):
		h.logger.Warn("send timeout", zap.String("conn_id", connID))
	}
}
