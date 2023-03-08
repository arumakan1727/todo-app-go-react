package restapi

import (
	"github.com/arumakan1727/todo-app-go-react/domain"
	oapi "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

func fillRespUser(r *RespUser, u *domain.User) {
	r.CreatedAt = u.CreatedAt
	r.DisplayName = u.DisplayName
	r.Email = oapi.Email(u.Email)
	r.Id = u.ID
}

func (h *UserHandler) ListUsersForAdmin(c echo.Context, clientUID UserID) error {
	ctx := c.Request().Context()

	xs, err := h.usecase.List(ctx)
	if err != nil {
		return err
	}
	return c.JSON(200, RespUserList{
		Items:      toRespArray(xs, fillRespUser),
		TotalCount: len(xs),
	})
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var b ReqCreateUser
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	u, err := h.usecase.Store(ctx, string(b.Email), b.Password, b.DisplayName, "user")
	if err != nil {
		return err
	}
	return c.JSON(200, toResp(&u, fillRespUser))
}
