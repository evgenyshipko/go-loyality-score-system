package server

import (
	"errors"
	c "github.com/evgenyshipko/go-loyality-score-system/internal/const"
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"github.com/evgenyshipko/go-loyality-score-system/internal/storage"
	"github.com/evgenyshipko/go-loyality-score-system/internal/tokens"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"net/http"
)

func (s *CustomServer) HelloWordHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.Write([]byte("<h1>Hello World</h1>"))
}

func (s *CustomServer) RegisterHandler(res http.ResponseWriter, req *http.Request) {

	data, err := getCredentialsFromContext(req.Context())
	if err != nil {
		logger.Instance.Warnw("getCredentialsFromContext", "err", err.Error())
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userId, err := s.services.Auth.Register(data.Login, data.Password)
	if err != nil {
		logger.Instance.Warnw("Ошибка Auth.Register", "err", err.Error())

		var pgErr pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			http.Error(res, "Пользователь с таким логином уже есть в базе", http.StatusConflict)
			return
		}

		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	s.registerTokens(userId, res)
}

func (s *CustomServer) LoginHandler(res http.ResponseWriter, req *http.Request) {
	data, err := getCredentialsFromContext(req.Context())
	if err != nil {
		logger.Instance.Warnw("getCredentialsFromContext", "err", err.Error())
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	success, userId, err := s.services.Auth.Login(data.Login, data.Password)
	if err != nil {
		logger.Instance.Warnw("Ошибка Auth.Login", "err", err.Error())

		var userNotFoundErr *storage.UserNotFoundError
		if errors.As(err, &userNotFoundErr) {
			http.Error(res, "Пользователя с таким логином нет в базе", http.StatusUnauthorized)
			return
		}

		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(res, "Неправильное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}
	s.registerTokens(userId, res)
}

func (s *CustomServer) LogoutHandler(res http.ResponseWriter, req *http.Request) {
	userId := req.Context().Value(c.UserId).(string)
	err := s.storage.DropUserTokens(userId)
	if err != nil {
		logger.Instance.Warnw("storage.DropUserTokens", "err", err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	clearCookies(res, []c.CookieName{c.AccessToken, c.RefreshToken})
}

func (s *CustomServer) RefreshHandler(res http.ResponseWriter, req *http.Request) {

	cookie, err := req.Cookie(string(c.RefreshToken))
	if err != nil {
		http.Error(res, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	claims, err := tokens.ParseJWT(cookie.Value)
	if err != nil {
		http.Error(res, "Неверный токен", http.StatusUnauthorized)
		return
	}

	s.registerTokens(claims.UserID, res)
}
