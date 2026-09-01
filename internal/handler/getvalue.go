package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h Handler) HandleGet(writer http.ResponseWriter, request *http.Request) {
	t := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")

	if t != "gauge" && t != "counter" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if t == "gauge" {
		handleGetGauge(h, writer, name)
	} else {
		handleGetCounter(h, writer, name)
	}
}

func handleGetCounter(h Handler, writer http.ResponseWriter, name string) {
	val, ok := h.storage.GetCounter(name)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write([]byte(val))
}

func handleGetGauge(h Handler, writer http.ResponseWriter, name string) {
	val, ok := h.storage.GetGauge(name)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write([]byte(val))
}
