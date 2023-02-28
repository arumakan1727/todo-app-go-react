package rest

import (
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
)

type AuthTokenHandler gHandler[domain.AuthTokenUsecase]

func (h AuthTokenHandler) CreateAuthToken(c echo.Context) error {
	ctx := c.Request().Context()

	var b domain.ReqCreateAuthToken
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	tok, err := h.usecase.Issue(ctx, &b)
	if err != nil {
		return err
	}
	return c.JSON(200, tok)
}
