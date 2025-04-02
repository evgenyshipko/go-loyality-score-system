package middlewares

import (
	"context"
	c "github.com/evgenyshipko/go-loyality-score-system/internal/const"
	"github.com/evgenyshipko/go-loyality-score-system/internal/tokens"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(string(c.AccessToken))
		if err != nil {
			http.Error(res, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		claims, err := tokens.ParseJWT(cookie.Value)
		if err != nil {
			http.Error(res, "Неверный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), c.UserId, claims.UserID)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
