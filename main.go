package main

import (
	"withjewel/jewel"
	"withjewel/handler"
)

func init() {
	jewel.Router("/demo", &handler.DemoHandler{})
}


func main() {
	jewel.Run(":8080")
}
