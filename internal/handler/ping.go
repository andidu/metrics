package handler

import (
	"net/http"
)

func (h Handler) HandlePing(writer http.ResponseWriter, request *http.Request) {
	isDBConnected := h.service.storage.Ping()

	if isDBConnected {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
