package handler

import (
	"withjewel/jewel"
)

/*IndexRequestHandler 处理主页请求 */
type IndexHandler struct {
	jewel.Controller
}

func (this *IndexHandler) Get() {
	this.RenderTpl("views/index.html", nil)
}
