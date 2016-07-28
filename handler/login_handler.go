package handler

import "withjewel/jewel"

/*LoginRequestHandler 处理登录请求 */
type LoginRequestHandler struct {
	jewel.Controller
}

/*Get 发送登录页面 */
func (this *LoginRequestHandler) Get() {
	username, _ := this.Cookie("username")
	password, _ := this.Cookie("password")

	this.Data["username"] = username
	this.Data["password"] = password
	//fmt.Println(this.Params)
	this.RenderTpl("views/login.html")
}

/*Post 验证登录*/
func (this *LoginRequestHandler) Post() {
	username := this.Ctx.Input.FormValue("username")
	password := this.Ctx.Input.FormValue("password")
	this.SetCookie("username", username)
	this.SetCookie("password", password)

	this.RenderTpl("views/login.html")
}
