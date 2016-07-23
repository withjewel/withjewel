package jewel

import (
	"net/http"
	"withjewel/jewel/jedb"
)

func Run(addr string) {
	jedb.Init()
	http.ListenAndServe(addr, nil)
}
