package jerrors

import (
    "errors"
)

var (
    QueryParamNotFound = errors.New("Query Param Not Found")
    CookieNotFound = errors.New("Cookie Not Found")
)
