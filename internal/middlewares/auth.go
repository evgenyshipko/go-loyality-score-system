package middlewares

import (
	"context"
	"github.com/evgenyshipko/go-loyality-score-system/internal/tokens"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("access_token")
		if err != nil {
			http.Error(res, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		claims, err := tokens.ParseJWT(cookie.Value)
		if err != nil {
			http.Error(res, "Неверный токен", http.StatusUnauthorized)
			return
		}

		// Передаем userID в контекст
		ctx := context.WithValue(req.Context(), "userID", claims.UserID)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
