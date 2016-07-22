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
    var datamodel = make(map[string]string)

    username, err := this.Cookie("username")
    fmt.Println(username)
    fmt.Println(err)
    if err == nil {
        datamodel["loginStatus"] = "true"
        datamodel["loginUser"] = username
    }

    this.RenderTpl("views/index.html", datamodel)
}
