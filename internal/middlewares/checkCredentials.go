package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"github.com/go-playground/validator"
	"net/http"
)

func CheckCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data c.Credentials
		err := validateBody(r, &data)
		if err != nil {
			logger.Instance.Warnw(err.Error())
			http.Error(w, "Не передан логин или пароль", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), c.CredentialsKey, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
