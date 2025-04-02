package server

import (
	"github.com/evgenyshipko/go-loyality-score-system/internal/db"
	"github.com/evgenyshipko/go-loyality-score-system/internal/httpserver"
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"github.com/evgenyshipko/go-loyality-score-system/internal/middlewares/logging"
	"github.com/evgenyshipko/go-loyality-score-system/internal/services"
	"github.com/evgenyshipko/go-loyality-score-system/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CustomServer struct {
	server   *httpserver.HTTPServer
	storage  *storage.SQLStorage
	services *services.Services
}

func NewCustomServer(router *chi.Mux) *CustomServer {

	// TODO: вынести serverDSN в переменные окружения
	database, err := db.ConnectToDB("postgres://metrics:metrics@localhost:5433/metrics?sslmode=disable&connect_timeout=5", true)
	if err != nil {
		panic(err)
	}

	store := storage.NewSQLStorage(database)
	service := services.NewServices(store)

	s := &CustomServer{
		// TODO: вынести host в переменные окружения
		server:   httpserver.NewHTTPServer("localhost:8080", router),
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
