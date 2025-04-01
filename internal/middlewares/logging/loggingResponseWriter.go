package logging

import (
	"bytes"
	"net/http"
)

type (
	// берём структуру для хранения сведений об ответе
	ResponseData struct {
		Status int
		Size   int
		Error  string // Здесь будет текст ошибки
	}

	// добавляем реализацию http.ResponseWriter
	LoggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		ResponseData        *ResponseData
		Body                bytes.Buffer // Буфер для тела ответа (используем только для ошибок)
		shouldCapture       bool         // Флаг, нужно ли сохранять тело
	}
)

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size // захватываем размер

	// Записываем тело только если статус ≥ 400
	if r.shouldCapture {
		r.Body.Write(b)
	}

	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)

	r.ResponseData.Status = statusCode // захватываем код статуса
	// Если статус ≥ 400, начинаем сохранять тело
	if statusCode >= 400 {
		r.shouldCapture = true
	}
}
