package main

import (
	"withjewel/handler"
	"withjewel/jewel"
    "fmt"
)

func init() {
	jewel.ServeStatic("static")

    jewel.Router("/", &handler.IndexHandler{})
	//jewel.Router("/user/<username:id>/page<pid:int>", &handler.IndexHandler{})
	//jewel.Router("/index", &handler.IndexHandler{})
	//jewel.Router("/login", &handler.LoginRequestHandler{})
	//jewel.Router("/logout", &handler.LogoutRequestHandler{})
}

func main() {
    fmt.Println("run")
	jewel.Run(":8080")
}
