package jerrors

import (
	"errors"
)

var (
	QueryStringNotFound = errors.New("Query Param Not Found")
	CookieNotFound      = errors.New("Cookie Not Found")
)
