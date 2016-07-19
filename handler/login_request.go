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
	fmt.Printf("用户请求登录，从%s\n", this.Input.RequestURI)
}
