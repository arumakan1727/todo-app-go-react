package domain

import "github.com/arumakan1727/todo-app-go-react/optional"

type TaskID uint64

type TaskPatch struct {
	Title optional.Option[string]
	Done  optional.Option[bool]
}

type TaskListFilter struct {
	DoneEq optional.Option[bool]
}

type TaskUcase interface {
	Store(ctx Ctx, uid UserID, title string) (Task, error)
	List(Ctx, UserID, TaskListFilter) ([]Task, error)
	Get(Ctx, UserID, TaskID) (Task, error)
	Patch(Ctx, UserID, TaskID, TaskPatch) (Task, error)
	Delete(Ctx, UserID, TaskID) error
}
