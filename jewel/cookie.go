package jewel

import (
	"errors"
	"net/http"
	"time"
)

const (
	// CookieDefaultMaxAge to one day.
	CookieDefaultMaxAge = 24 * 60 * 60
)

var (
	CookieDefaultExpires = time.Now().Add(CookieDefaultMaxAge * time.Second)
)

var (
	// Cookie errors.
	CookieExpiredError = errors.New("Expired Cookie Error")
)

// NewCookie to make a new cookie.
func (controller *Controller) NewCookie(name string, value string, others ...interface{}) (cookie *http.Cookie) {
	cookie = &http.Cookie{
		Name:  name,
		Value: value,

		Path:    "/",
		MaxAge:  CookieDefaultMaxAge,
		Expires: CookieDefaultExpires,
	}
	return cookie
}

// SetCookie to set one cookie for the response.
func (controller *Controller) SetCookie(name string, value string, others ...interface{}) {
	cookie := controller.NewCookie(name, value, others...)
	http.SetCookie(controller.Ctx.Output, cookie)
}

// CookieObject to get the *cookie with the specific name.
func (controller *Controller) CookieObject(name string) (cookie *http.Cookie, err error) {
	cookie, err = controller.Ctx.Input.Cookie(name)
	if err != nil {
		return nil, err
	}
	// if cookie.MaxAge <= 0 {
	// 	return nil, CookieExpiredError
	// }
	return cookie, err
}

// Cookie to get the value of the specific cookie.
func (controller *Controller) Cookie(name string) (value string, err error) {
	cookie, err := controller.CookieObject(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// RemCookie to remove a cookie.
func (controller *Controller) RemCookie(name string) {
	cookie := controller.NewCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(controller.Ctx.Output, cookie)
}
