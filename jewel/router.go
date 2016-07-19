package jewel

import (
	"net/http"
)

type Handler struct {
	ControllerInterface
	Input  *http.Request
	Output http.ResponseWriter
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.Init(rw, r)
	switch {
	case r.Method == "GET":
		h.Get()
	case r.Method == "POST":
		h.Post()
	}
}

func Router(url string, controllerInterface ControllerInterface) {
	handler := &Handler{controllerInterface, nil, nil}
	http.Handle(url, handler)
}
