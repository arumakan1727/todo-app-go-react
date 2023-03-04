package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/labstack/echo/v4"
)

func parseBodyAsJSON(ctx context.Context, r *http.Request, dest interface{}) error {
	if r.Body == nil || r.Body == http.NoBody {
		return Err400EmptyReqBody
	}
	defer r.Body.Close()

	// CORS 対策のために application/json であることを要請する。
	// 例えば text/plain の POST だと単純リクエスト扱いになってしまい preflight が送信されない。
	// 'charset=utf-8' を許容するために HasPrefix() を使う。
	ct := r.Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(ct, echo.MIMEApplicationJSON) {
		return Err400InvalidReqContentType
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(dest); err != nil {
		return Err400UndecodableJSON
	}
	return nil
}

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
