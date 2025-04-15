package server

import (
	_config "github.com/evgenyshipko/go-rag-chat-helper/internal/config"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/db"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/httpserver"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/middlewares/logging"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/services"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CustomServer struct {
	server   *httpserver.HTTPServer
	storage  *storage.SQLStorage
	services *services.Services
}

func NewCustomServer(router *chi.Mux) *CustomServer {
	config := _config.GetConfig()

	database, err := db.ConnectToDB(config.PostgresConnect)
	if err != nil {
		panic(err)
	}

	store := storage.NewSQLStorage(database)
	service := services.NewServices(store, &config)

	s := &CustomServer{
		server:   httpserver.NewHTTPServer(config.ServerHost, router),
		storage:  store,
		services: service,
	}

	apiRouter := s.initRoutes()

	router.Mount("/api", apiRouter)

	return s
}

func Create() *CustomServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Use(logging.LoggingHandlers)

	server := NewCustomServer(router)
	return server
}

func (s *CustomServer) Start() {
	err := s.server.Start()
	if err != nil {
		logger.Instance.Warn("Failed to start server")
	}
}

func (s *CustomServer) ShutDown() {
	logger.Instance.Info("Завершение сервера...")

	err := s.server.Stop()
	if err != nil {
		logger.Instance.Warnw("CustomServer.Shutdown", "Ошибка завершения сервера Stop()", err)
	}

	logger.Instance.Info("Сервер успешно завершён")
}
