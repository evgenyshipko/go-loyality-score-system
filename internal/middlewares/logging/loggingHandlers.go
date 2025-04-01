package logging

import (
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"github.com/evgenyshipko/go-loyality-score-system/internal/middlewares/utils"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func LoggingHandlers(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Получаем requestID из контекста (middleware chi)
		requestID := middleware.GetReqID(r.Context())

		requestBody, err := utils.GetBodyAndRestore(r)
		if err != nil {
			logger.Instance.Warnw("GetBodyAndRestore", "err", err)
			http.Error(w, "GetBodyAndRestore ошибка", http.StatusBadRequest)
			return
		}

		loggerFunc := logger.Instance.Infow

		loggerFunc("Request",
			"requestId", requestID,
			"uri", r.RequestURI,
			"method", r.Method,
			"body", requestBody,
		)

		responseData := &ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := LoggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			ResponseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		if responseData.Status >= 400 {
			loggerFunc = logger.Instance.Warnw
		}

		loggerFunc("Response",
			"requestId", requestID,
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
			"body", lw.Body.String(),
		)
	}
	return http.HandlerFunc(logFn)
}
