package jewel

import (
	"fmt"
	"net/http"
	"net/url"
	"withjewel/jewel/jerrors"
)

type Controller struct {
	Name   string
	Ctx    RequestContext

	Params url.Values
}

type RequestContext struct {
	Input *http.Request
	Output http.ResponseWriter
}

type ControllerInterface interface {
	Init(http.ResponseWriter, *http.Request)
	Get()
	Post()
}

func (this *Controller) Init(responseWriter http.ResponseWriter, request *http.Request) {
	this.Ctx.Input = request
	this.Ctx.Output = responseWriter

	this.Params = request.URL.Query()
}

/*Cookie 相关*/

// 添加cookie
func (this *Controller) SetCookie(cookie *http.Cookie) {
	if cookie.MaxAge <= 0 {
		cookie.MaxAge = 99999999
		fmt.Println(cookie.MaxAge)
	}
	http.SetCookie(this.Ctx.Output, cookie)
}

// 取出cookie的值
func (this *Controller) Cookie(name string) (string, error) {
	cookie, err := this.CookieObject(name)
	if err != nil || cookie.MaxAge <= 0 {
		return "", jerrors.CookieNotFound
	}
	return cookie.Value, nil
}

// 取出cookie指针
func (this *Controller) CookieObject(name string) (*http.Cookie, error) {
	cookie, err := this.Ctx.Input.Cookie(name)
	return cookie, err
}

// 移除cookie
func (this *Controller) RemoveCookie(name string) {
	cookie, err := this.CookieObject(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	// MaxAge < 0 means remove the cookie right now
	cookie.MaxAge = -200
}


/*Param 相关*/

// 取出param的值
func (this *Controller) ParamString(name string) (string, error) {
	if param, ok := this.Params[name]; ok {
		return param[0], nil
	}
	return "", jerrors.QueryParamNotFound
}

/*模板处理*/

// 渲染文件模板
func (this *Controller) RenderTpl(tplPath string, model interface{}) {
	RenderTplFile(this.Ctx.Output, tplPath, model)
}

/*Redirect*/

func (this *Controller) Redirect(url string) {
	http.Redirect(this.Ctx.Output, this.Ctx.Input, url, 302)
}



// 默认的GET请求处理函数
func (h *Controller) Get() {
	http.Error(h.Ctx.Output, "Method Not Allowed", 405)
}

// 默认的POST请求处理函数
func (h *Controller) Post() {
	http.Error(h.Ctx.Output, "Method Not Allowed", 405)
}
