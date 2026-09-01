package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h Handler) HandleUpdate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")
	strvalue := chi.URLParam(request, "value")

	if t != "gauge" && t != "counter" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if t == "gauge" {
		handleUpdateGauge(h, writer, name, strvalue)
	} else {
		handleUpdateCounter(h, writer, name, strvalue)
	}
}

func handleUpdateGauge(h Handler, writer http.ResponseWriter, name, strvalue string) {
	value, error := strconv.ParseFloat(strvalue, 64)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.UpdateGauge(name, value)
	writer.WriteHeader(http.StatusOK)
}

func handleUpdateCounter(h Handler, writer http.ResponseWriter, name, strvalue string) {
	value, error := strconv.Atoi(strvalue)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	h.storage.UpdateCounter(name, value)
	writer.WriteHeader(http.StatusOK)
}
