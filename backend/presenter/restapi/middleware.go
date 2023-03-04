package restapi

import (
	"net/http"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
)

type MiddlewareFunc = func(next echo.HandlerFunc) echo.HandlerFunc

func AuthMiddleware(au domain.AuthUsecase) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			cookie, err := c.Cookie(cookieKeyAuthToken)
			if err != nil {
				return c.String(http.StatusUnauthorized, "Missing auth token.")
			}

			uid, err := au.ValidateAuthToken(ctx, domain.AuthToken(cookie.Value))
			if err != nil {
				return c.String(http.StatusUnauthorized, "Invalid auth token.")
			}

			ctxWithUID := newCtxWithUserID(ctx, uid)
			c.SetRequest(c.Request().WithContext(ctxWithUID))
			return next(c)
		}
	}
}
