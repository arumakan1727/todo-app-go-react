package restapi

import (
	"net/http"
	"time"

	"github.com/arumakan1727/todo-app-go-react/config"
)

const CookieKeyAuthToken = "todoApiAuthToken"

func newSecureCookie(
	name, value string, maxAge time.Duration, runMode config.RunMode,
) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   "",
		MaxAge:   int(maxAge.Seconds()),
		Secure:   runMode != config.ModeDebug,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func deleteCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
