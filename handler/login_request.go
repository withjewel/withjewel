package handler

import (
	"fmt"
	"regexp"
	"withjewel/jewel"
	"withjewel/model"
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
	UsernameRegex := regexp.MustCompile(`[A-Za-z0-9]{4}`)
	EmailRegex := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

	username := this.Ctx.Input.FormValue("username")
	password := this.Ctx.Input.FormValue("password")

	datamodel := make(map[string]string)

	var pass = false

	if username == "" {
		datamodel["loginStatus"] = "失败"
		datamodel["loginMessage"] = "你小子肯定是通过非法客户端连接的。"
	} else {
		if len(UsernameRegex.FindString(username)) == len(username) {
			pass = model.Verify(username, password)
			if !pass {
				datamodel["loginStatus"] = "失败"
				datamodel["loginMessage"] = "错误的用户名或密码，或许您需要填写一个电子邮件来加入？"
			}
		} else if len(EmailRegex.FindString(username)) == len(username) {
			pass = model.VerifyEmail(username, password)
			if !pass {
				datamodel["loginStatus"] = "成功"
				datamodel["loginMessage"] = "我们已经向" + username + "发送了一封确认电子邮件。"
			}
		} else {
			datamodel["loginStatus"] = "失败"
			datamodel["loginMessage"] = "你小子肯定是通过非法客户端连接的。"
		}
	}

	if pass {
		datamodel["loginStatus"] = "成功"
		datamodel["loginMessage"] = "欢迎您"
	}

	//fmt.Printf("%s请求登录, 密码%s\n", username, password)
	jewel.RenderTplFile(this.Ctx.Output, "views/login.html", datamodel)
}
