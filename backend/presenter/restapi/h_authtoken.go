package restapi

import (
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
)

type AuthTokenHandler struct {
	runMode config.RunMode
	usecase domain.AuthUsecase
}

func (h *AuthTokenHandler) CreateAuthToken(c echo.Context) error {
	ctx := c.Request().Context()

	var b ReqCreateAuthToken
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	token, err := h.usecase.IssueAuthToken(ctx, string(b.Email), b.Password)
	if err != nil {
		return err
	}
	c.SetCookie(newSecureCookie(
		cookieKeyAuthToken, string(token),
		h.usecase.GetAuthTokenMaxAge(),
		h.runMode,
	))
	return c.String(200, "OK")
}
