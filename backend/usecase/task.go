package usecase

import (
	. "github.com/arumakan1727/todo-app-go-react/domain"
)

type taskUc struct {
}

func NewTaskUsecase() TaskUcase {
	return &taskUc{}
}

func (uc *taskUc) Store(ctx Ctx, uid UserID, title string) (Task, error) {
	panic("TODO")
}

func (uc *taskUc) List(ctx Ctx, uid UserID, f TaskListFilter) ([]Task, error) {
	panic("TODO")
}

func (uc *taskUc) Get(Ctx, UserID, TaskID) (Task, error) {
	panic("TODO")
}

func (uc *taskUc) Patch(Ctx, UserID, TaskID, TaskPatch) (Task, error) {
	panic("TODO")
}

func (uc *taskUc) Delete(Ctx, UserID, TaskID) error {
	panic("TODO")
}
