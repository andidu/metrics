package handler

import (
	"fmt"
	"net/http"
)

func (h Handler) HandleDefault(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	gauges, err := h.service.Gauges()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	counters, err := h.service.Counters()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write([]byte("<h2>gauge</h2>"))
	for name, value := range gauges {
		writer.Write(fmt.Appendf([]byte{}, "<div>%s - %f</div>", name, value))
	}

	writer.Write([]byte("<h2>counter</h2>"))
	for name, value := range counters {
		writer.Write(fmt.Appendf([]byte{}, "<div>%s - %d</div>", name, value))
	}
}
