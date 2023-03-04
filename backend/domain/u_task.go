package domain

import (
	"time"

	"github.com/arumakan1727/todo-app-go-react/optional"
)

type TaskID uint64

type TaskPatch struct {
	Title optional.Option[string]
	Done  optional.Option[bool]
}

type TaskListFilter struct {
	DoneEq optional.Option[bool]
}

type TaskUsecase interface {
	Store(ctx Ctx, uid UserID, title string) (Task, error)
	List(Ctx, UserID, TaskListFilter) ([]Task, error)
	Get(Ctx, UserID, TaskID) (Task, error)
	Patch(Ctx, UserID, TaskID, TaskPatch) (Task, error)
	Delete(Ctx, UserID, TaskID) (Task, error)
}

func (t *Task) ApplyTimezone(loc *time.Location) {
	t.CreatedAt = t.CreatedAt.In(loc)
}
