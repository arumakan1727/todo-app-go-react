package restapi

import (
	"net/http"
	"strings"

	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
)

type MiddlewareFunc = func(next echo.HandlerFunc) echo.HandlerFunc

func CORSMiddleware(cfg *config.Config) MiddlewareFunc {
	setAccessControlHeader := func(h http.Header, origin string) {
		h.Set(echo.HeaderAccessControlAllowOrigin, origin)
		h.Set(echo.HeaderAccessControlAllowMethods, "GET,HEAD,POST,PUT,PATCH,DELETE")
		h.Set(echo.HeaderAccessControlAllowHeaders, strings.Join([]string{
			echo.HeaderAcceptEncoding,
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderCookie,
		}, ","))
		h.Set(echo.HeaderAccessControlMaxAge, "86400") // 86400sec = 24h
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Debug("called CORSMiddleware")
			origin := c.Request().Header.Get(echo.HeaderOrigin)
			if len(origin) == 0 {
				return next(c)
			}

			isOptionsMethod := c.Request().Method == http.MethodOptions

			// 127.0.0.1 もしくは localhost からのリクエストは許可
			addr := strings.TrimPrefix(origin, "http://")
			if strings.HasPrefix(addr, "127.0.0.1") || strings.HasPrefix(addr, "localhost") {
				if !cfg.AllowLocalhostOrigin {
					return c.NoContent(http.StatusForbidden)
				}
				setAccessControlHeader(c.Response().Header(), origin)
				if isOptionsMethod {
					return c.NoContent(200)
				}
				return next(c)
			}

			for _, o := range cfg.AllowedOrigins {
				if origin == o {
					setAccessControlHeader(c.Response().Header(), origin)
					if isOptionsMethod {
						return c.NoContent(200)
					}
					return next(c)
				}
			}
			c.Logger().Info("Request from unallowed origin:", origin)
			return c.NoContent(http.StatusForbidden)
		}
	}
}

func AuthMiddleware(au domain.AuthUsecase) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Debug("called AuthMiddleware")
			ctx := c.Request().Context()

			cookie, err := c.Cookie(cookieKeyAuthToken)
			if err != nil {
				// return c.String(http.StatusUnauthorized, "Missing auth token.")
				return next(c)
			}

			am, err := au.ValidateAuthToken(ctx, domain.AuthToken(cookie.Value))
			if err != nil {
				// return c.String(http.StatusUnauthorized, "Invalid auth token.")
				return next(c)
			}

			ctxWithUID := newCtxWithAuthMaterial(ctx, am)
			c.SetRequest(c.Request().WithContext(ctxWithUID))
			return next(c)
		}
	}
}

func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		am, err := getAuthMaterialFromCtx(c.Request().Context())
		if err != nil || !am.IsAdmin() {
			// 存在自体を知られないようにするために Not Found
			return echo.NewHTTPError(http.StatusNotFound, "Not Found")
		}
		return next(c)
	}
}
