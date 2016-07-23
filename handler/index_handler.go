package handler

import (
    "fmt"
    "withjewel/jewel"
)

/*IndexRequestHandler 处理主页请求 */
type IndexHandler struct {
    jewel.Controller
}

func (this *IndexHandler) Get() {
    fmt.Println("username=", this.ParamString("username"))
    fmt.Println("page_id=", this.ParamInt("pid"))
    var datamodel = make(map[string]string)

    username, err := this.Cookie("username")
    if err == nil {
        datamodel["loginStatus"] = "true"
        datamodel["loginUser"] = username
    }

    this.RenderTpl("views/index.html", datamodel)
}
