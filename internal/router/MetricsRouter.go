package router

import (
	"github.com/andidu/metrics/internal/handler"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func MetricsRouter(handler handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	r.Get("/value/{type}/{name}", handler.HandleGet)
	r.HandleFunc("/update/{type}/{name}/{value}", handler.HandleUpdate)
	r.HandleFunc("/", handler.HandleDefault)

	return r
}
