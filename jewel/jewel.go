package jewel

import (
	"log"
	"net/http"
	"strings"
	"withjewel/jewel/jedb"
)

const (
	defaultHost = "localhost"
	addrSplit   = ":"
)

func init() {
	jedb.Init()
}

func parseListenAddr(addr string) (host string, port string) {
	addrFields := strings.Split(addr, addrSplit)
	host, port = addrFields[0], addrFields[1]
	if host == "" {
		host = defaultHost
	}
	return
}

// Run 开启HTTP服务器
func Run(addr string) {
	host, port := parseListenAddr(addr)
	log.Println("[server]: Start listening...")
	log.Println("[server]:", host)
	log.Println("[server]:", port)
	http.ListenAndServe(addr, DefaultJewelServeMux)
}
