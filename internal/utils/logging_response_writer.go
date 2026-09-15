package utils

import "net/http"

type LoggingResponseWriter struct {
	http.ResponseWriter
	data *ResponseData
}

type ResponseData struct {
	status int
	size   int
}

func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.data.status = statusCode
}
