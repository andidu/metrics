package main

import (
	"net/http"

	"github.com/andidu/metrics/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/counter/", handler.HandleUpdateCounter)
	mux.HandleFunc("/update/gauge/", handler.HandleUpdateGauge)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		println("Server didn't start")
	}
}
