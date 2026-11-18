package hash

import (
	"net/http"
	"slices"
)

type HashWriter struct {
	http.ResponseWriter
	BuffHeader     http.Header
	BuffStatusCode int
	BuffResponse   []byte
}

func (hw *HashWriter) Write(b []byte) (int, error) {
	hw.BuffResponse = slices.Concat(hw.BuffResponse, b)
	return len(b), nil
}

func (hw *HashWriter) WriteHeader(statusCode int) {
	hw.BuffStatusCode = statusCode
}

func (hw *HashWriter) Header() http.Header {
	return hw.BuffHeader
}
