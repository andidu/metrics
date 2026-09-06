package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h Handler) HandleUpdate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	t := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")
	strvalue := chi.URLParam(request, "value")

	err := h.service.UpdateMetric(t, name, strvalue)
	switch err {
	case UnknownMetricTypeErr:
		writer.WriteHeader(http.StatusBadRequest)
		return
	case WrongMetricValueErr:
		writer.WriteHeader(http.StatusBadRequest)
		return
	case InternalStorageErr:
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	writer.WriteHeader(http.StatusOK)
}
