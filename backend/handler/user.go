package handler

import (
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
)

type UserHandler Handler[domain.UserUsecase]

func (h UserHandler) ListUsersForAdmin(c echo.Context) error {
	ctx := c.Request().Context()

	ul, err := h.usecase.List(ctx)
	if err != nil {
		return err
	}
	return c.JSON(200, ul)
}

func (h UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var b domain.ReqCreateUser
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	user, err := h.usecase.Store(ctx, &b)
	if err != nil {
		return err
	}
	return c.JSON(200, user)
}
