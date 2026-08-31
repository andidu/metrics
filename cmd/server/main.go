package main

import (
	"net/http"

	"github.com/andidu/metrics/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/counter/{name}", handler.HandleGetCounter)
	r.Get("/gauge/{name}", handler.HandleGetGauge)
	r.HandleFunc("/update/counter/{name}/{value}", handler.HandleUpdateCounter)
	r.HandleFunc("/update/gauge/{name}/{value}", handler.HandleUpdateGauge)
	r.HandleFunc("/", handler.HandleDefault)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		println("Server didn't start", err.Error())
	}
}
