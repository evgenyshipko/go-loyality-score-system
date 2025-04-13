package server

import (
	"bytes"
	"errors"
	"fmt"
	c "github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/storage"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/tokens"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"io"
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
		http.Error(res, "Отсутствует рефреш-токен", http.StatusUnauthorized)
		return
	}

	claims, err := tokens.ParseJWT(cookie.Value)
	if err != nil {
		http.Error(res, "Неверный токен", http.StatusUnauthorized)
		return
	}

	s.registerTokens(claims.UserID, res)
}

func (s *CustomServer) UploadHandler(w http.ResponseWriter, r *http.Request) {
	//Ограничиваем размер загружаемого файла до 10 МБ
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при ParseMultipartForm: %v", err), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при получении файла: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при чтении файла: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.services.Document.UploadDocument(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Чанки загружены успешно!"))
}

func (s *CustomServer) AnswerHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		logger.Instance.Warnw("chi.URLParam", "err", "query is empty")
		http.Error(w, "необходимо передать запрос в параметре query", http.StatusBadRequest)
		return
	}

	keywords, err := s.services.Llm.GetKeywords(query)
	if err != nil {
		logger.Instance.Warnw("Llm.GetKeywords", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	document, err := s.services.Document.SearchDocument(keywords)
	if document == "" {
		w.Write([]byte("Извините, я не знаю ответа на данный вопрос. Могу отвечать только на вопросы по компании."))
		return
	}

	if err != nil {
		logger.Instance.Warnw("searchDocument", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	answer, err := s.services.Llm.GetAnswerBasedOnDocument(query, document)
	if err != nil {
		logger.Instance.Warnw("llm.GetAnswerBasedOnDocument", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(answer))
}
