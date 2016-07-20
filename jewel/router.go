package jewel

import (
	"net/http"
)

type Handler struct {
	ControllerInterface
}

func (h *Handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.Init(responseWriter, request)
	switch {
	case request.Method == "GET":
		h.Get()
	case request.Method == "POST":
		h.Post()
	}
}

func Router(url string, controllerInterface ControllerInterface) {
	handler := &Handler{controllerInterface}
	http.Handle(url, handler)
}
