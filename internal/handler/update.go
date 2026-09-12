package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	models "github.com/andidu/metrics/internal/model"
)

func (h Handler) HandleUpdate(writer http.ResponseWriter, request *http.Request) {
	var metrics models.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metrics)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateMetrics(metrics)
	createUpdateResponse(err, writer)
}
