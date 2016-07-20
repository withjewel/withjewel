package handler

import (
	"fmt"
	"withjewel/jewel"
)

/*LoginRequestHandler 处理登录请求 */
type LoginRequestHandler struct {
	jewel.Controller
}

/*Get 发送登录页面 */
func (this *LoginRequestHandler) Get() {
	jewel.RenderTplFile(this.Ctx.Output, "views/login.html", nil)
	fmt.Printf("用户请求登录，从%s\n", this.Ctx.Input.RequestURI)
}

/*Post 验证登录*/
func (this *LoginRequestHandler) Post() {
	username := this.Input.FormValue("username")
	fmt.Printf("%s请求登录\n", username)
	this.Output.Write(nil)
}
