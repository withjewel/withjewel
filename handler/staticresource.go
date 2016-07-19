package handler

import (
    "withjewel/jewel"
    "net/http"
    "fmt"
)

type StaticResourceHandler struct {
    jewel.Controller
}

func (this *StaticResourceHandler) Get() {
    fmt.Printf("请求文件%s\n", this.Input.URL.Path)
    bytes, err := jewel.Get(this.Input.URL.Path)
    if err != nil {
        http.Error(this.Output, "Resource not found", 404)
    }
    this.Output.Write(bytes)
}