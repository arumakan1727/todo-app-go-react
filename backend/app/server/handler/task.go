package handler

import (
	"context"

	"github.com/arumakan1727/todo-app-go-react/app/server/schema"
	"github.com/labstack/echo/v4"
)

type TaskUsecase interface {
	List(ctx context.Context, params schema.ListTasksParams) (*schema.TaskList, error)
	Create(ctx context.Context, title string) (*schema.Task, error)
	Get(ctx context.Context, tid schema.TaskID) (*schema.Task, error)
	Patch(ctx context.Context, tid schema.TaskID, title *string, done *bool) (*schema.Task, error)
	Delete(ctx context.Context, tid schema.TaskID) error
}

type TaskHandler Handler[TaskUsecase]

func (h *TaskHandler) ListTasks(c echo.Context, params schema.ListTasksParams) error {
	ctx := c.Request().Context()

	list, err := h.usecase.List(ctx, params)
	if err != nil {
		return err
	}
	return c.JSON(200, list)
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	ctx := c.Request().Context()

	var b schema.ReqCreateTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	task, err := h.usecase.Create(ctx, b.Title)
	if err != nil {
		return err
	}
	return c.JSON(200, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context, taskID schema.TaskID) error {
	ctx := c.Request().Context()

	err := h.usecase.Delete(ctx, taskID)
	if err != nil {
		return err
	}
	return c.String(200, "deleted")
}

func (h *TaskHandler) GetTask(c echo.Context, taskID schema.TaskID) error {
	ctx := c.Request().Context()

	task, err := h.usecase.Get(ctx, taskID)
	if err != nil {
		return err
	}
	return c.JSON(200, task)
}

func (h *TaskHandler) PatchTask(c echo.Context, taskID schema.TaskID) error {
	ctx := c.Request().Context()

	var b schema.ReqPatchTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	task, err := h.usecase.Patch(ctx, taskID, b.Title, b.Done)
	if err != nil {
		return err
	}
	return c.JSON(200, task)
}
