package handler

import (
	"context"

	"github.com/arumakan1727/todo-app-go-react/app/server/schema"
	"github.com/labstack/echo/v4"
)

type AuthTokenUsecase interface {
	Create(ctx context.Context, email, passwd string) (*schema.AuthToken, error)
}

type AuthTokenHandler Handler[AuthTokenUsecase]

func (h AuthTokenHandler) CreateAuthToken(c echo.Context) error {
	ctx := c.Request().Context()

	var b schema.ReqCreateAuthToken
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	tok, err := h.usecase.Create(ctx, string(b.Email), b.Password)
	if err != nil {
		return err
	}
	return c.JSON(200, tok)
}
