package jewel

import (
	"withjewel/jewel/jerrors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)




type ControllerInterface interface {
	Init(http.ResponseWriter, *http.Request)
	InitParams(map[string]string)
	Get()
	Post()
}

type RequestContext struct {
	Input *http.Request
	Output http.ResponseWriter
}

type Controller struct {
	Name   string
	Ctx    RequestContext

	QueryString url.Values
	Params      map[string]string
}


func (this *Controller) Init(responseWriter http.ResponseWriter, request *http.Request, ) {
	this.Ctx.Input = request
	this.Ctx.Output = responseWriter

	this.QueryString = request.URL.Query()
}

func (this *Controller) InitParams(params map[string]string) {
	this.Params = params
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

/*Query字符串 相关*/

// 取出param的值
func (this *Controller) Query(name string) (string, error) {
	if queryString, ok := this.QueryString[name]; ok {
		return queryString[0], nil
	}
	return "", jerrors.QueryStringNotFound
}

/*Url参数 相关*/

func (this *Controller) ParamString(name string) string {
	if paramStr, ok := this.Params[name]; ok {
		return paramStr
	}
	return ""
}

func (this *Controller) ParamInt(name string) int {
	if paramStr, ok := this.Params[name]; ok {
		value, err := strconv.Atoi(paramStr)
		if err == nil {
			return value
		}
	}
	return 0
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
