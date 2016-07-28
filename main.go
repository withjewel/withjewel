package main

import (
	"withjewel/handler"
	"withjewel/jewel"
)

func init() {
	jewel.ServeStatic("/static", "static")

	//jewel.Router("/login", &handler.LoginRequestHandler{})
	jewel.Router("/logout", &handler.LogoutRequestHandler{})
	jewel.Router("/user", &handler.UserHandler{})
	jewel.Router("/", &handler.IndexHandler{})
}

func main() {
	jewel.Run(":8080")
}
