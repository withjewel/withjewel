package jewel

import (
	"net/http"
)

type Controller struct {
	Name   string
	Ctx    RequestContext
}

type RequestContext struct {
	Input *http.Request
	Output http.ResponseWriter
}

type ControllerInterface interface {
	Init(http.ResponseWriter, *http.Request)
	Get()
	Post()
}

func (h *Controller) Init(responseWriter http.ResponseWriter, request *http.Request) {
	h.Ctx.Input = request
	h.Ctx.Output = responseWriter
}

func (h *Controller) Get() {
	http.Error(h.Ctx.Output, "Method Not Allowed", 405)
}

func (h *Controller) Post() {
	http.Error(h.Ctx.Output, "Method Not Allowed", 405)
}
