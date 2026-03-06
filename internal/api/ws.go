package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"Philograph/internal/application"

	"nhooyr.io/websocket"
)

// WSHub はWebSocket接続を管理しブロードキャストする。
type WSHub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]context.CancelFunc
}

// NewWSHub は新しいWSHubを返す。
func NewWSHub() *WSHub {
	return &WSHub{
		clients: make(map[*websocket.Conn]context.CancelFunc),
	}
}

// HandleWS はWebSocket接続を処理するHTTPハンドラー。
func (h *WSHub) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("websocket accept error", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())

	h.mu.Lock()
	h.clients[conn] = cancel
	h.mu.Unlock()

	slog.Info("websocket client connected", "clients", h.clientCount())

	// Keep connection alive, read messages (for ping/pong)
	go func() {
		defer func() {
			h.removeClient(conn)
			conn.Close(websocket.StatusNormalClosure, "")
		}()
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				return
			}
		}
	}()
}

// Broadcast は全クライアントに進捗情報を送信する。
func (h *WSHub) Broadcast(progress application.Progress) {
	data, err := json.Marshal(progress)
	if err != nil {
		slog.Error("failed to marshal progress", "error", err)
		return
	}

	h.mu.RLock()
	clients := make(map[*websocket.Conn]context.CancelFunc, len(h.clients))
	for c, cancel := range h.clients {
		clients[c] = cancel
	}
	h.mu.RUnlock()

	for conn := range clients {
		ctx, cancel := context.WithCancel(context.Background())
		err := conn.Write(ctx, websocket.MessageText, data)
		cancel()
		if err != nil {
			h.removeClient(conn)
		}
	}
}

// ProgressListener は ProgressListener 型のコールバックを返す。
func (h *WSHub) ProgressListener() application.ProgressListener {
	return func(p application.Progress) {
		h.Broadcast(p)
	}
}

func (h *WSHub) removeClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if cancel, ok := h.clients[conn]; ok {
		cancel()
		delete(h.clients, conn)
	}
	slog.Info("websocket client disconnected", "clients", len(h.clients))
}

func (h *WSHub) clientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
