package handler

import (
	"fmt"
	"time"
	"withjewel/jewel"
)

type DemoHandler struct {
	jewel.Controller
}

func (d *DemoHandler) Get() {
	fmt.Println("Get /demo")
	model := make(map[string]interface{})
	model["WebTitle"] = "with jewel"
	model["Date"] = time.Now()
	model["Gays"] = []string{"Jiaju.chen", "Nan.li"}
	jewel.RenderTplFile(d.Ctx.Output, "views/demo.html", model)
}
