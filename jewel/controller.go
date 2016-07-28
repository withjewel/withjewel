package jewel

import (
	"net/http"
	"net/url"
	"strconv"
	"withjewel/jewel/jerrors"
)

type Context struct {
	Input  *http.Request
	Output http.ResponseWriter
}

type Controller struct {
	Name string
	Ctx  Context

	QueryString url.Values
	Params      map[string]string
	Data        map[string]interface{}
}

func (this *Controller) Init(responseWriter http.ResponseWriter, request *http.Request) {
	this.Ctx.Input = request
	this.Ctx.Output = responseWriter

	this.QueryString = request.URL.Query()
	this.Data = make(map[string]interface{})
}

func (this *Controller) InitParams(params map[string]string) {
	this.Params = params
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
func (this *Controller) RenderTpl(tplPath string) {

	RenderTplFile(this.Ctx.Output, tplPath, this.Data)
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

// 默认的PUT请求处理函数
func (h *Controller) Put() {
	http.Error(h.Ctx.Output, "Method Not Allowed", 405)
}

// 默认的PATCH请求处理函数
func (h *Controller) Patch() {
	http.Error(h.Ctx.Output, "Method Not Allowed", 405)
}

// 默认的DELETE请求处理函数
func (h *Controller) Delete() {
	http.Error(h.Ctx.Output, "Method Not Allowed", 405)
}
