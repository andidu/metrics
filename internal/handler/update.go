package handler

import (
	"net/http"
	"strconv"

	"github.com/andidu/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

func HandleUpdate(writer http.ResponseWriter, request *http.Request) {
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
		handleUpdateGauge(writer, name, strvalue)
	} else {
		handleUpdateCounter(writer, name, strvalue)
	}
}

func handleUpdateGauge(writer http.ResponseWriter, name, strvalue string) {
	value, error := strconv.ParseFloat(strvalue, 64)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	service.TmpInMemoryStarage.UpdateGauge(name, value)
	writer.WriteHeader(http.StatusOK)
}

func handleUpdateCounter(writer http.ResponseWriter, name, strvalue string) {
	value, error := strconv.Atoi(strvalue)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	service.TmpInMemoryStarage.UpdateCounter(name, value)
	writer.WriteHeader(http.StatusOK)
}
