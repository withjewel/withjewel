package handler

import (
	"withjewel/jewel"
)

type UserHandler struct {
	jewel.Controller
}

func (this *UserHandler) Get() {
	username, _ := this.Cookie("username")
	password, _ := this.Cookie("password")

	this.Data["username"] = username
	this.Data["password"] = password
	this.RenderTpl("views/user.html")
}
