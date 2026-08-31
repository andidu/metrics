package main

import (
	"net/http"

	"github.com/andidu/metrics/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{type}/{name}", handler.HandleGet)
	r.HandleFunc("/update/{type}/{name}/{value}", handler.HandleUpdate)
	r.HandleFunc("/", handler.HandleDefault)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		println("Server didn't start", err.Error())
	}
}
