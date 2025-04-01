package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"github.com/evgenyshipko/go-loyality-score-system/internal/storage"
	"github.com/go-playground/validator"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"net/http"
	"time"
)

func (s *CustomServer) HelloWordHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.Write([]byte("<h1>Hello World</h1>"))
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func validateBody[T any](req *http.Request, data *T) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf.Bytes(), data); err != nil {
		return err
	}

	validate := validator.New() // создаём валидатор

	if err := validate.Struct(data); err != nil {
		return errors.New(fmt.Sprintf("Не заполнены поля: %s", err))
	}

	return nil
}

func (s *CustomServer) RegisterHandler(res http.ResponseWriter, req *http.Request) {

	var data RegisterRequest
	err := validateBody(req, &data)
	if err != nil {
		logger.Instance.Warnw(err.Error())
		http.Error(res, "Не передан логин или пароль", http.StatusBadRequest)
		return
	}

	err = s.services.Auth.Register(data.Login, data.Password)
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
	s.registerTokens(data.Login, res)
}

func (s *CustomServer) LoginHandler(res http.ResponseWriter, req *http.Request) {
	var data RegisterRequest
	err := validateBody(req, &data)
	if err != nil {
		logger.Instance.Warnw(err.Error())
		http.Error(res, "Не передан логин или пароль", http.StatusBadRequest)
		return
	}

	success, err := s.services.Auth.Login(data.Login, data.Password)
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
	s.registerTokens(data.Login, res)
}

func (s *CustomServer) registerTokens(login string, res http.ResponseWriter) {
	access, refresh, err := s.services.Auth.GenerateTokensAndSave(login)
	if err != nil {
		logger.Instance.Warnw("Ошибка Auth.GenerateTokensAndSave", "err", err.Error())

		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	http.SetCookie(res, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
}
