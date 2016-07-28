package jewel

import (
	"net/http"
)

const (
	DefaultAddr = "localhost:8080"
)

var JewelAPP = NewApp()

type Jewel struct {
	*http.Server
}

// NewApp return a new jewel app.
func NewApp() *Jewel {
	server := new(http.Server)
	server.Addr = DefaultAddr
	server.Handler = DefaultJewelServeMux

	app := &Jewel{server}
	return app
}

// Run Jewel app.
func Run(addr string) {
	JewelAPP.Addr = addr
	JewelAPP.ListenAndServe()
}
