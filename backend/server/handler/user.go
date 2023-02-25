package handler

import (
	"context"

	"github.com/arumakan1727/todo-app-go-react/schema"
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	ListAll(c context.Context) (*schema.UserList, error)
	Create(c context.Context, param *schema.ReqCreateUser) (*schema.User, error)
}

type UserHandler Handler[UserUsecase]

func (h UserHandler) ListUsersForAdmin(c echo.Context) error {
	ctx := c.Request().Context()

	ul, err := h.usecase.ListAll(ctx)
	if err != nil {
		return err
	}
	return c.JSON(200, ul)
}

func (h UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var b schema.ReqCreateUser
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	user, err := h.usecase.Create(ctx, &b)
	if err != nil {
		return err
	}
	return c.JSON(200, user)
}
