package jewel

import (
	"fmt"
	"net/http"
	"strings"
)

//StaticFileHandler 静态文件处理
type StaticFileHandler struct {
	Controller
	DirName   string
	URIPrifix string
}

// Get 处理对于静态文件的GET请求
func (h *StaticFileHandler) Get() {
	filename := strings.TrimLeft(h.Ctx.Input.RequestURI, "/")
	http.ServeFile(h.Ctx.Output, h.Ctx.Input, filename)
}

// ServeStatic 注册静态文件路由和对应的静态文件文件夹
func ServeStatic(pattern string, dirName string) {
	handler := new(StaticFileHandler)
	handler.DirName = dirName
	handler.URIPrifix = fmt.Sprintf("/%s/", dirName)

	Router(pattern, handler)
}
