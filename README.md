# withjewel
网站 http://www.withjewel.com 的golang后端

### TODO
* 实现对于cookie和session的基本支持
* 实现对于Url pattern的正则支持

### 说明
暂时没有使用任何第三方web框架, 之后会根据具体要求不断完善jewel这个自建简易“框架”吧。

### 简易使用方法

在handler文件夹下添加对于具体url和方法的处理模块，然后在main.go中注册该url和该handler的对应关系。
大致过程如下:

**新增handler模块**

```go

// hello_world.go
package handler

import (
    "withjewel/jewel"
)

type DemoHandler struct {
    jewel.Controller
}

func (d *DemoHandler) Get() {
    fmt.Println("Get /demo")
    model := make(map[string]interface{})
    model["WebTitle"] = "with jewel"
    fmt.Println(model)
    //d.Output.Write([]byte("helo"))
    jewel.RenderTplFile(d.Output, "views/hello.html", model)
}
```

**新建html模板文件**

```html

<html>
    <head>
        <title>{{.WebTitle}}</title>
    </head>
    <body>
        <h1>Welcome to {{.WebTitle}}</h1>
    </body>
</html>
```

**在main.go中注册路由***

```go
func init() {
    jewel.Router("/demo", &handler.DemoHandler{})
}
```
