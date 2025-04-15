package constants

type ContextKey string

const CredentialsKey ContextKey = "credentials"
const UserId ContextKey = "userId"

type Credentials struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CookieName string

const AccessToken CookieName = "access_token"
const RefreshToken CookieName = "refresh_token"

type LlmRoles string

const User LlmRoles = "user"
const System LlmRoles = "system"

type DocumentChunk struct {
	Text       string
	ChunkIndex int
}
