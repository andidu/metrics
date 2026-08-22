package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/andidu/metrics/internal/service"
)

func HandleUpdateGauge(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(request.URL.Path, "/update/gauge/")
	data := strings.Split(path, "/")
	if len(data) != 2 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	name := data[0]
	if len(name) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	value, error := strconv.ParseFloat(data[1], 64)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	service.TmpInMemoryStarage.UpdateGauge(name, value)
	writer.WriteHeader(http.StatusOK)
}

func HandleUpdateCounter(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(request.URL.Path, "/update/counter/")
	data := strings.Split(path, "/")
	if len(data) != 2 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	name := data[0]
	if len(name) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	value, error := strconv.Atoi(data[1])
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	service.TmpInMemoryStarage.UpdateCounter(name, value)
	writer.WriteHeader(http.StatusOK)
}
