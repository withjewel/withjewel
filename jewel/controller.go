package jewel

import (
	"net/http"
)

type Controller struct {
	Name   string
	Output http.ResponseWriter
	Input  *http.Request
}

type ControllerInterface interface {
	Init(http.ResponseWriter, *http.Request)
	Get()
	Post()
}

func (h *Controller) Init(rw http.ResponseWriter, req *http.Request) {
	h.Output = rw
	h.Input = req
}

func (h *Controller) Get() {
	http.Error(h.Output, "Method Not Allowed", 405)
}

func (h *Controller) Post() {
	http.Error(h.Output, "Method Not Allowed", 405)
}
