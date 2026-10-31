package handler

import (
	"net/http"
)

func (h Handler) HandlePing(writer http.ResponseWriter, request *http.Request) {
	isDbConnected := h.service.storage.Ping()

	if isDbConnected {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
