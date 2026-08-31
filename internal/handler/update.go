package handler

import (
	"net/http"
	"strconv"

	"github.com/andidu/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

func HandleUpdateGauge(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain")

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := chi.URLParam(request, "name")
	strvalue := chi.URLParam(request, "value")

	value, error := strconv.ParseFloat(strvalue, 64)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	service.TmpInMemoryStarage.UpdateGauge(name, value)
	writer.WriteHeader(http.StatusOK)
}

func HandleUpdateCounter(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain")

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := chi.URLParam(request, "name")
	strvalue := chi.URLParam(request, "value")

	value, error := strconv.Atoi(strvalue)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	service.TmpInMemoryStarage.UpdateCounter(name, value)
	writer.WriteHeader(http.StatusOK)
}
