package handler

type Handler struct {
	storage MemStorage
}

func New(s MemStorage) Handler {
	return Handler{
		storage: s,
	}
}
