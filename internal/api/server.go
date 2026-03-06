package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"Philograph/web"
)

// Server はHTTPサーバー。
type Server struct {
	handler *Handler
	wsHub   *WSHub
	server  *http.Server
	addr    string
}

// NewServer は新しいServerを返す。
func NewServer(handler *Handler, wsHub *WSHub, port int) *Server {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/v1/analyze", handler.HandleAnalyze)
	mux.HandleFunc("PATCH /api/v1/config", handler.HandleUpdateConfig)
	mux.HandleFunc("GET /api/v1/config", handler.HandleGetConfig)
	mux.HandleFunc("GET /api/v1/result", handler.HandleGetResult)
	mux.HandleFunc("GET /api/v1/export", handler.HandleExport)
	mux.HandleFunc("GET /api/v1/health", handler.HandleHealth)

	// WebSocket
	mux.HandleFunc("GET /ws", wsHub.HandleWS)

	// Static files (SPA fallback)
	mux.Handle("/", web.StaticHandler())

	// Apply middleware
	var h http.Handler = mux
	h = loggingMiddleware(h)
	h = recoveryMiddleware(h)

	addr := fmt.Sprintf(":%d", port)

	return &Server{
		handler: handler,
		wsHub:   wsHub,
		addr:    addr,
		server: &http.Server{
			Addr:    addr,
			Handler: h,
		},
	}
}

// Start はサーバーを開始する。ポート0の場合は自動割り当て。
func (s *Server) Start() (int, error) {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return 0, err
	}

	port := ln.Addr().(*net.TCPAddr).Port
	slog.Info("server starting", "port", port)

	go func() {
		if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	return port, nil
}

// Shutdown はサーバーをグレースフルにシャットダウンする。
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("server shutting down")
	return s.server.Shutdown(ctx)
}
