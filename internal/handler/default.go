package handler

import (
	"fmt"
	"net/http"

	"github.com/andidu/metrics/internal/service"
)

func HandleDefault(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusBadRequest)
	}
	writer.Write([]byte("<h2>gauge</h2>"))
	for name, value := range service.TmpInMemoryStarage.Gauges {
		writer.Write(fmt.Appendf([]byte{}, "<div>%s - %f</div>", name, value))
	}

	writer.Write([]byte("<h2>counter</h2>"))
	for name, value := range service.TmpInMemoryStarage.Counters {
		writer.Write(fmt.Appendf([]byte{}, "<div>%s - %d</div>", name, value))
	}
}
