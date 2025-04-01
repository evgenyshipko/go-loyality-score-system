package middlewares

import (
	"compress/gzip"
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"net/http"
)

func GzipDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Instance.Warn(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close()
		}

		next.ServeHTTP(w, r)
	})
}
