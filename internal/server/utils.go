package server

import (
	"context"
	"errors"
	c "github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"net/http"
	"time"
)

func (s *CustomServer) registerTokens(userId string, res http.ResponseWriter) {
	access, refresh, err := s.services.Auth.GenerateTokensAndSave(userId)
	if err != nil {
		logger.Instance.Warnw("Ошибка Auth.GenerateTokensAndSave", "err", err.Error())

		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     string(c.AccessToken),
		Value:    access,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	http.SetCookie(res, &http.Cookie{
		Name:     string(c.RefreshToken),
		Value:    refresh,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
}

func clearCookies(w http.ResponseWriter, names []c.CookieName) {
	for _, cookieName := range names {
		http.SetCookie(w, &http.Cookie{
			Name:    string(cookieName),
			Value:   "",
			Path:    "/",
			MaxAge:  -1,
			Expires: time.Unix(0, 0),
		})
	}
}

func getCredentialsFromContext(ctx context.Context) (c.Credentials, error) {
	metricData := ctx.Value(c.CredentialsKey)

	data, ok := metricData.(c.Credentials)
	if !ok {
		logger.Instance.Warn("Невозможно привести к Credentials")
		return c.Credentials{}, errors.New("невозможно привести к Credentials")
	}
	return data, nil
}
