package gzip

import (
	"io"
	"net/http"
)

type GzipWriter struct {
	http.ResponseWriter
	W io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.W.Write(b)
}
