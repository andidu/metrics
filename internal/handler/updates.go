package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	models "github.com/andidu/metrics/internal/model"
)

func (h Handler) HandleUpdates(writer http.ResponseWriter, request *http.Request) {
	var metrics = make([]models.Metrics, 0)
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

	err = h.service.UpdateMetricsBatch(metrics)
	createUpdateResponse(err, writer)
}
