package main

import (
	"net/http"

	config "github.com/andidu/metrics/internal/config/server"
	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/router"
	"github.com/andidu/metrics/internal/service"
	"github.com/andidu/metrics/internal/utils"
)

func main() {
	utils.InitLogger()
	defer utils.Logger.Sync()

	var flags = config.ParseConfig()

	storage := service.NewMemStorage()
	handler := handler.New(storage)
	r := router.MetricsRouter(handler)

	err := http.ListenAndServe(flags.ServerAddress, utils.WithLogging(r))
	if err != nil {
		println("Server didn't start", err.Error())
	}
}
