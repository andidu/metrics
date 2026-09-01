package main

import (
	"net/http"

	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/router"
	"github.com/andidu/metrics/internal/service"
	"github.com/spf13/pflag"
)

func main() {
	pflag.Parse()

	storage := service.NewMemStorage()
	handler := handler.New(storage)
	r := router.MetricsRouter(handler)

	err := http.ListenAndServe(*serverAddress, r)
	if err != nil {
		println("Server didn't start", err.Error())
	}
}
