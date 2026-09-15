package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	models "github.com/andidu/metrics/internal/model"
)

func (h Handler) Value(writer http.ResponseWriter, request *http.Request) {
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

	err = h.service.RetreiveMetrics(&metrics)
	if err != nil {
		switch err {
		case ErrUnknownMetricType:
			writer.WriteHeader(http.StatusBadRequest)
			return
		case ErrWrongMetricValue:
			writer.WriteHeader(http.StatusBadRequest)
			return
		case ErrInternalStorage:
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		case ErrNoElement:
			writer.WriteHeader(http.StatusNotFound)
			return
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
	}

	resp, err := json.Marshal(metrics)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(resp)
}
