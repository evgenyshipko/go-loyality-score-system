package utils

import (
	"bytes"
	"io"
	"net/http"
)

func GetBodyAndRestore(r *http.Request) (string, error) {
	var requestBody string
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		requestBody = string(bodyBytes)
		// Восстанавливаем r.Body, чтобы обработчики могли его использовать
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	return requestBody, nil
}
