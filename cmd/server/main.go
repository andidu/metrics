package main

import (
	"net/http"

	"github.com/andidu/metrics/internal/router"
)

func main() {
	r := router.MetricsRouter()

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		println("Server didn't start", err.Error())
	}
}
