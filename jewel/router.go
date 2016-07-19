package jewel

import (
	"net/http"
	"fmt"
)

type Handler struct {
	ControllerInterface
	Input *http.Request
	Output http.ResponseWriter
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.Init(rw, r)
	switch r.Method {
	case "GET":
		h.Get()
	case "POST":
		h.Post()
	}
}

func Router(url string, controllerInterface ControllerInterface) {
	handler := &Handler{controllerInterface, nil, nil}
	fmt.Printf("url: %s, handler: %s\n", url, handler)
	http.Handle(url, handler)
}
