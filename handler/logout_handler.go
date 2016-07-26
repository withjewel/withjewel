package handler

import "withjewel/jewel"

/*LogoutRequestHandler 处理登出请求 */
type LogoutRequestHandler struct {
	jewel.Controller
}

func (this *LogoutRequestHandler) Get() {
	this.RemCookie("username")
	this.RemCookie("password")
	this.Redirect("/index")
}
