package rest

import (
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
)

type TaskHandler gHandler[domain.TaskUcase]

func (h TaskHandler) ListTasks(c echo.Context, params domain.ListTasksParams) error {
	ctx := c.Request().Context()

	list, err := h.usecase.List(ctx, params)
	if err != nil {
		return err
	}

	res := make([]domain.RespTask, 0, len(list))
	for _, t := range list {
		res = append(res, domain.RespTask{
			CreatedAt: time.Time{},
			Done:      false,
			Id:        0,
			Title:     "",
		})
	}
	return c.JSON(200, list)
}

func (h TaskHandler) CreateTask(c echo.Context) error {
	ctx := c.Request().Context()

	var b domain.ReqCreateTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	task, err := h.usecase.Store(ctx, &b)
	if err != nil {
		return err
	}
	return c.JSON(200, task)
}

func (h TaskHandler) DeleteTask(c echo.Context, taskID domain.TaskID) error {
	ctx := c.Request().Context()

	err := h.usecase.Delete(ctx, 0, taskID)
	if err != nil {
		return err
	}
	return c.String(200, "deleted")
}

func (h TaskHandler) GetTask(c echo.Context, taskID domain.TaskID) error {
	ctx := c.Request().Context()

	task, err := h.usecase.Get(ctx, 0, taskID)
	if err != nil {
		return err
	}
	return c.JSON(200, task)
}

func (h TaskHandler) PatchTask(c echo.Context, taskID domain.TaskID) error {
	ctx := c.Request().Context()

	var b domain.ReqPatchTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	task, err := h.usecase.Patch(ctx, 0, taskID, &b)
	if err != nil {
		return err
	}
	return c.JSON(200, task)
}
