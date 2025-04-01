package services

import "github.com/evgenyshipko/go-loyality-score-system/internal/storage"

type Services struct {
	Auth *AuthService
}

func NewServices(storage *storage.SQLStorage) *Services {
	return &Services{
		Auth: NewAuthService(storage),
	}
}
