package server

import (
	"github.com/evgenyshipko/go-loyality-score-system/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func (s *CustomServer) initRoutes() *chi.Mux {
	apiRouter := chi.NewRouter()

	apiRouter.With(middlewares.AuthMiddleware).Get("/", s.HelloWordHandler)

	apiRouter.Post("/user/register", s.RegisterHandler)

	apiRouter.Post("/user/login", s.LoginHandler)

	return apiRouter
}
