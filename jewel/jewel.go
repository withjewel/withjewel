package jewel

import (
	"net/http"
)

func Run(addr string) {
	http.ListenAndServe(addr, nil)
}
