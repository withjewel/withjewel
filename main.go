package main

import (
	"withjewel/handler"
	"withjewel/jewel"
)

func init() {
	jewel.Router("/demo", &handler.DemoHandler{})
	jewel.Router("/login", &handler.LoginRequestHandler{})
}

func main() {
	jewel.Run(":8080")
}
