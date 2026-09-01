package handler

import (
	"github.com/andidu/metrics/internal/service"
)

type Handler struct {
	storage service.MemStorage
}

func New(s service.MemStorage) Handler {
	return Handler{
		storage: s,
	}
}
