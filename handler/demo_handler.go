package handler

import (
	"../../withjewel/jewel"
	"time"
	"fmt"
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
	fmt.Println(model)
	//d.Output.Write([]byte("helo"))
	jewel.RenderTplFile(d.Output, "views/demo.html", model)
}
