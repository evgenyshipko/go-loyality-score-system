package httpserver

import (
	"context"
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"net/http"
	"time"
)

type HTTPServer struct {
	server http.Server
}

func NewHTTPServer(host string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: http.Server{Addr: host, Handler: handler},
	}
}

func (s *HTTPServer) Start() error {
	logger.Instance.Infow("SERVER STARTED!")
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Instance.Fatalw("httpServer.ListenAndServe", "Ошибка запуска сервера", err)
		return err
	}
	return nil
}

func (s *HTTPServer) Stop() error {
	// Создаём контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Instance.Warnw("httpServer.Shutdown", "Ошибка завершения сервера:", err)
		return err
	}
	return nil
}
