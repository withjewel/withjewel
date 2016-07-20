package jewel

import (
	"fmt"
	"net/http"
	"strings"
)

type StaticFileHandler struct {
	Controller
	DirName   string
	URIPrifix string
}

func (h *StaticFileHandler) Get() {
	filename := strings.TrimLeft(h.Ctx.Input.RequestURI, "/")
	http.ServeFile(h.Ctx.Output, h.Ctx.Input, filename)
}

/*用于注册静态文件文件夹*/
func ServeStatic(dirName string) {
	handler := new(StaticFileHandler)
	handler.DirName = dirName
	handler.URIPrifix = fmt.Sprintf("/%s/", dirName)

	Router(handler.URIPrifix, handler)
}
