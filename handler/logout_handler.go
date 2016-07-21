package handler

import (
    "fmt"
    "withjewel/jewel"
)

/*LogoutRequestHandler 处理登出请求 */
type LogoutRequestHandler struct {
    jewel.Controller
}

func (this *LogoutRequestHandler) Get() {
    this.RemoveCookie("username")
    fmt.Println(this.Cookie("username"))
    this.Redirect("/index")
}
