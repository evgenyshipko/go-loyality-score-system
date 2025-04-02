package server

import (
	"github.com/evgenyshipko/go-loyality-score-system/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

func (s *CustomServer) initRoutes() *chi.Mux {
	apiRouter := chi.NewRouter()

	apiRouter.With(middlewares.Auth).Get("/", s.HelloWordHandler)

	apiRouter.With(middlewares.Auth).Post("/user/logout", s.LogoutHandler)

	apiRouter.With(middlewares.CheckCredentials).Post("/user/register", s.RegisterHandler)

	apiRouter.With(middlewares.CheckCredentials).Post("/user/login", s.LoginHandler)

	apiRouter.Post("/user/refresh", s.RefreshHandler)

	return apiRouter
}
