package handler

import "withjewel/jewel"

type DemoHandler struct {
	jewel.Controller
}

func (this *DemoHandler) Get() {
	this.RenderTpl("views/demo.html")
}
