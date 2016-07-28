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
	fmt.Println(this.Query("name"))
	this.Data["time"] = time.Now()
	time.Sleep(5 * time.Second)
	fmt.Println(this.Data)
	this.RenderTpl("views/index.html")
}
