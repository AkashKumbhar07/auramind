package market

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var defaultDialer = &websocket.Dialer{
	HandshakeTimeout: 10 * time.Second,
}

type WebSocketConn struct {
	conn   *websocket.Conn
	mu     sync.Mutex
	logger interface {
		Info(string, ...any)
		Error(string, ...any)
	}
}

func Dial(_ context.Context, url string) (*WebSocketConn, error) {
	header := http.Header{}
	conn, _, err := defaultDialer.Dial(url, header)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", url, err)
	}
	return &WebSocketConn{conn: conn}, nil
}

func (w *WebSocketConn) Read() ([]byte, error) {
	_, msg, err := w.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	return msg, nil
}

func (w *WebSocketConn) Write(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *WebSocketConn) WriteJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return w.Write(data)
}

func (w *WebSocketConn) Close() error {
	return w.conn.Close()
}

func (w *WebSocketConn) SetLogger(l interface {
	Info(string, ...any)
	Error(string, ...any)
}) {
	w.logger = l
}

func jsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
