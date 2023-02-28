package rest

import (
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
)

type TaskHandler gHandler[domain.TaskUcase]

func fillRespTask(r *RespTask, t *domain.Task) {
	r.CreatedAt = t.CreatedAt
	r.Done = t.Done
	r.Id = t.ID
	r.Title = t.Title
}

func (h TaskHandler) ListTasks(c echo.Context, params ListTasksParams) error {
	ctx := c.Request().Context()
	uid, err := ctxGetUserID(ctx)
	if err != nil {
		return err
	}

	var filterDoneEq *bool
	if params.Status == nil || *params.Status == TaskStatusFilterAny {
		filterDoneEq = nil
	} else {
		filterDoneEq = new(bool)
		*filterDoneEq = (*params.Status == TaskStatusFilterDone)
	}
	xs, err := h.usecase.List(ctx, uid, filterDoneEq)
	if err != nil {
		return err
	}

	return c.JSON(200, RespTaskList{
		Items:      toRespArray(xs, fillRespTask),
		TotalCount: len(xs),
	})
}

func (h TaskHandler) CreateTask(c echo.Context) error {
	ctx := c.Request().Context()
	uid, err := ctxGetUserID(ctx)
	if err != nil {
		return err
	}

	var b ReqCreateTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	task, err := h.usecase.Store(ctx, uid, b.Title)
	if err != nil {
		return err
	}
	return c.JSON(200, toResp(&task, fillRespTask))
}

func (h TaskHandler) DeleteTask(c echo.Context, taskID domain.TaskID) error {
	ctx := c.Request().Context()
	uid, err := ctxGetUserID(ctx)
	if err != nil {
		return err
	}

	err = h.usecase.Delete(ctx, uid, taskID)
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
	return c.JSON(200, toResp(&task, fillRespTask))
}

func (h TaskHandler) PatchTask(c echo.Context, taskID domain.TaskID) error {
	ctx := c.Request().Context()
	uid, err := ctxGetUserID(ctx)
	if err != nil {
		return err
	}

	var b ReqPatchTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	patch := domain.TaskPatch{
		Title: b.Title,
		Done:  b.Done,
	}
	task, err := h.usecase.Patch(ctx, uid, taskID, patch)
	if err != nil {
		return err
	}
	return c.JSON(200, toResp(&task, fillRespTask))
}
