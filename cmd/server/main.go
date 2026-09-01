package main

import (
	"net/http"

	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/router"
	"github.com/andidu/metrics/internal/service"
)

func main() {
	storage := service.NewMemStorage()
	handler := handler.New(storage)
	r := router.MetricsRouter(handler)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		println("Server didn't start", err.Error())
	}
}
