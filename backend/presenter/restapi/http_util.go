package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
		if emsg := err.Error(); strings.HasPrefix(emsg, "email:") {
			return echo.NewHTTPError(400, emsg)
		}
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			emsg := fmt.Sprintf("%s: unexpected value type (got: %s, offset: %d)", e.Field, e.Value, e.Offset)
			return echo.NewHTTPError(400, emsg)
		}
		if e, ok := err.(*json.SyntaxError); ok {
			emsg := fmt.Sprintf("%s (offset=%d)", e.Error(), e.Offset)
			return echo.NewHTTPError(400, emsg)
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return echo.NewHTTPError(400, err.Error())
		}
		log.Printf("parseBodyAsJSON(): unknown err: %#v", err)
		return echo.NewHTTPError(400, "invalid request JSON")
	}
	return nil
}
