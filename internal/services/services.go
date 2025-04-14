package services

import "github.com/evgenyshipko/go-rag-chat-helper/internal/storage"

type Services struct {
	Auth     *AuthService
	Document *DocumentService
	Llm      *LlmService
}

func NewServices(storage *storage.SQLStorage) *Services {
	return &Services{
		Auth:     NewAuthService(storage),
		Document: NewDocumentService(storage),
		Llm:      NewLlmService(),
	}
}
