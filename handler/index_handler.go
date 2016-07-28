package handler

import (
	"fmt"
	"time"
	"withjewel/jewel"
)

/*IndexRequestHandler 处理主页请求 */
type IndexHandler struct {
	jewel.Controller
}

func (this *IndexHandler) Get() {
	this.Data["time"] = time.Now()
	fmt.Println(this.Data)
	this.RenderTpl("views/index.html")
}
