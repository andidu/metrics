package handler

import (
	"net/http"

	"github.com/andidu/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

func HandleGetCounter(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain")

	metric := chi.URLParam(request, "name")

	val, ok := service.TmpInMemoryStarage.GetCounter(metric)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write([]byte(val))
}

func HandleGetGauge(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain")

	metric := chi.URLParam(request, "name")

	val, ok := service.TmpInMemoryStarage.GetGauge(metric)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write([]byte(val))
}
