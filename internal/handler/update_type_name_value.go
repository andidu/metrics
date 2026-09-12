package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h Handler) HandleUpdateTypeNameValue(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")
	strvalue := chi.URLParam(request, "value")

	err := h.service.UpdateMetric(t, name, strvalue)
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
	}

	writer.WriteHeader(http.StatusOK)
}
