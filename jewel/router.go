package jewel

import (
	"net/http"
)

// JewelHandler

type JewelHandler struct {
	ControllerInterface
}

func (h *JewelHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.Init(responseWriter, request)
	switch {
	case request.Method == "GET":
		h.Get()
	case request.Method == "POST":
		h.Post()
	}
}

func Router(pattern string, controllerInterface ControllerInterface) {
	handler := &JewelHandler{controllerInterface}
	DefaultJewelServeMux.Handle(pattern, handler)
}
