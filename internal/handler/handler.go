package handler

type Handler struct {
	service Service
}

func New(s MemStorage) Handler {
	return Handler{
		service: Service{
			storage: s,
		},
	}
}
