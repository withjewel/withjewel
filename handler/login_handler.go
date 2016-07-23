package handler

import (
	"fmt"
	"net/http"
	"withjewel/jewel"
	"withjewel/model"
	"withjewel/utils"
)

/*LoginRequestHandler 处理登录请求 */
type LoginRequestHandler struct {
	jewel.Controller
}

/*Get 发送登录页面 */
func (this *LoginRequestHandler) Get() {
	datamodel := make(map[string]string)
	username, err := this.Cookie("username")

	if err == nil && username == "leslie" {
		datamodel["loginStatus"] = "true"
		datamodel["loginUser"] = username
		fmt.Println(datamodel)
	} else {
		jewel.RenderTplFile(this.Ctx.Output, "views/login.html", nil)
	}
}

/*Post 验证登录*/
func (this *LoginRequestHandler) Post() {
	username := this.Ctx.Input.FormValue("username")
	password := this.Ctx.Input.FormValue("password")

	datamodel := make(map[string]string)

	var pass = false
	if utils.VerifyUsername(username) {
		pass = model.Verify(username, password)
		if !pass {
			datamodel["loginStatus"] = "失败"
			datamodel["loginMessage"] = "错误的用户名或密码，或许您需要填写一个电子邮件来加入？"
		}
	} else if utils.VerifyEmail(username) {
		exist, pass := model.VerifyEmail(username, password)
		if !exist {
			datamodel["loginStatus"] = "成功"
			datamodel["loginMessage"] = "我们已经向" + username + "发送了一封确认电子邮件。"
		} else if !pass {
			datamodel["loginStatus"] = "失败"
			datamodel["loginMessage"] = "错误的密码。"
		}
	} else {
		datamodel["loginStatus"] = "失败"
		datamodel["loginMessage"] = "你小子肯定是通过非法客户端连接的。"
	}

	if pass {
		datamodel["loginStatus"] = "成功"
		datamodel["loginMessage"] = "欢迎您"
		cookie := &http.Cookie{Name: "username", Value: username}
		this.SetCookie(cookie)
		//this.Redirect("/index")
		//return
	}

	//fmt.Printf("%s请求登录, 密码%s\n", username, password)
	jewel.RenderTplFile(this.Ctx.Output, "views/login.html", datamodel)
}
