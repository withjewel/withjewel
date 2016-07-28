package jewel

import (
	"net/http"
	"reflect"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodPatch  = "patch"
	MethodDelete = "Delete"
)

// Controller is a interface.
type ControllerInterface interface {
	Init(http.ResponseWriter, *http.Request)
	InitParams(map[string]string)
	Get()
	Post()
	Put()
	Patch()
	Delete()
}

type HandlerFactory struct {
	SeedType reflect.Type
}

func (hf *HandlerFactory) newControllerInterface() ControllerInterface {
	nc, _ := reflect.New(hf.SeedType).Interface().(ControllerInterface)
	return nc
}

func (hf *HandlerFactory) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	c := hf.newControllerInterface()
	c.Init(rw, request)
	switch request.Method {
	case MethodGet:
		c.Get()
	case MethodPost:
		c.Post()
	case MethodPut:
		c.Put()
	case MethodPatch:
		c.Patch()
	case MethodDelete:
		c.Delete()
	}
}

// Router register the pattern with a handler.
func Router(pattern string, c ControllerInterface) {
	reflectValue := reflect.ValueOf(c)
	seedType := reflect.Indirect(reflectValue).Type()

	handlerFactory := &HandlerFactory{SeedType: seedType}

	DefaultJewelServeMux.Handle(pattern, handlerFactory)
}
