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
	name, strvalue, parsed := parseParams(writer, path)
	if !parsed {
		return
	}

	value, error := strconv.ParseFloat(strvalue, 64)
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
	name, strvalue, parsed := parseParams(writer, path)
	if !parsed {
		return
	}

	value, error := strconv.Atoi(strvalue)
	if error != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	service.TmpInMemoryStarage.UpdateCounter(name, value)
	writer.WriteHeader(http.StatusOK)
}

// Returns name and value for a metric. The bool indecates whether the extraction happened correctly
func parseParams(writer http.ResponseWriter, path string) (string, string, bool) {
	data := strings.Split(path, "/")
	if len(data) > 2 {
		writer.WriteHeader(http.StatusNotFound)
		return "", "", false
	}
	if len(data) < 2 {
		writer.WriteHeader(http.StatusNotFound)
		return "", "", false
	}

	name := data[0]
	if len(name) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return "", "", false
	}

	return name, data[1], true
}
