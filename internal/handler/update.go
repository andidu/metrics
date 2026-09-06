package handler

import (
	"log"
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
	value, err := strconv.ParseFloat(strvalue, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.UpdateGauge(name, value)
	if err != nil {
		log.Println("Internal storage error", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func handleUpdateCounter(h Handler, writer http.ResponseWriter, name, strvalue string) {
	value, err := strconv.Atoi(strvalue)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.UpdateCounter(name, value)
	if err != nil {
		log.Println("Internal storage error", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	writer.WriteHeader(http.StatusOK)
}
