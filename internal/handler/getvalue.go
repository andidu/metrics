package handler

import (
	"net/http"

	"github.com/andidu/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

func HandleGet(writer http.ResponseWriter, request *http.Request) {
	t := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")

	if t != "gauge" && t != "counter" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if t == "gauge" {
		handleGetGauge(writer, name)
	} else {
		handleGetCounter(writer, name)
	}
}

func handleGetCounter(writer http.ResponseWriter, name string) {
	val, ok := service.TmpInMemoryStarage.GetCounter(name)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write([]byte(val))
}

func handleGetGauge(writer http.ResponseWriter, name string) {
	val, ok := service.TmpInMemoryStarage.GetGauge(name)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write([]byte(val))
}
