package services

import (
	"github.com/evgenyshipko/go-rag-chat-helper/internal/config"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/storage"
)

type ServiceNeedsStorage interface {
}

type Services struct {
	Auth     *AuthService
	Document *DocumentService
	Llm      *LlmService
}

func NewServices(storage *storage.SQLStorage, config *config.Config) *Services {
	return &Services{
		Auth:     NewAuthService(storage),
		Document: NewDocumentService(storage),
		Llm:      NewLlmService(config),
	}
}
