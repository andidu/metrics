package utils

import (
	"net/http"
	"time"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		uri := r.RequestURI
		method := r.Method

		responseData := &ResponseData{}

		lw := LoggingResponseWriter{
			ResponseWriter: w,
			data:           responseData,
		}

		h.ServeHTTP(&lw, r)

		totalTime := time.Since(startTime)

		Logger.Infoln(
			"uri", uri,
			"method", method,
			"time", totalTime,
			"status", responseData.status,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
