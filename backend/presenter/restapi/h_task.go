package restapi

import (
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/optional"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	usecase domain.TaskUcase
}

func fillRespTask(r *RespTask, t *domain.Task) {
	r.CreatedAt = t.CreatedAt
	r.Done = t.Done
	r.Id = t.ID
	r.Title = t.Title
}

func (h TaskHandler) ListTasks(
	c echo.Context, clientUID UserID, params ListTasksParams,
) error {
	ctx := c.Request().Context()

	var f domain.TaskListFilter
	if params.Status == nil || *params.Status == TaskStatusFilterAny {
		f.DoneEq = optional.None[bool]()
	} else {
		f.DoneEq = optional.Some(*params.Status == TaskStatusFilterDone)
	}
	xs, err := h.usecase.List(ctx, clientUID, f)
	if err != nil {
		return err
	}

	return c.JSON(200, RespTaskList{
		Items:      toRespArray(xs, fillRespTask),
		TotalCount: len(xs),
	})
}

func (h TaskHandler) CreateTask(c echo.Context, clientUID UserID) error {
	ctx := c.Request().Context()

	var b ReqCreateTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	task, err := h.usecase.Store(ctx, clientUID, b.Title)
	if err != nil {
		return err
	}
	return c.JSON(200, toResp(&task, fillRespTask))
}

func (h TaskHandler) DeleteTask(
	c echo.Context, clientUID UserID, taskID domain.TaskID,
) error {
	ctx := c.Request().Context()

	err := h.usecase.Delete(ctx, clientUID, taskID)
	if err != nil {
		return err
	}
	return c.String(200, "deleted")
}

func (h TaskHandler) GetTask(
	c echo.Context, clientUID UserID, taskID domain.TaskID,
) error {
	ctx := c.Request().Context()

	task, err := h.usecase.Get(ctx, clientUID, taskID)
	if err != nil {
		return err
	}
	return c.JSON(200, toResp(&task, fillRespTask))
}

func (h TaskHandler) PatchTask(
	c echo.Context, clientUID UserID, taskID domain.TaskID,
) error {
	ctx := c.Request().Context()

	var b ReqPatchTask
	if err := parseBodyAsJSON(ctx, c.Request(), &b); err != nil {
		return err
	}

	patch := domain.TaskPatch{
		Title: optional.FromNillable(b.Title),
		Done:  optional.FromNillable(b.Done),
	}
	task, err := h.usecase.Patch(ctx, clientUID, taskID, patch)
	if err != nil {
		return err
	}
	return c.JSON(200, toResp(&task, fillRespTask))
}
