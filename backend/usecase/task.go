package usecase

import (
	"context"

	. "github.com/arumakan1727/todo-app-go-react/domain"
)

type taskUcase struct {
}

func NewTaskUsecase() TaskUcase {
	return &taskUcase{}
}

func (uc *taskUcase) List(ctx context.Context, p ListTasksParams) (TaskList, error) {
	return TaskList{}, nil
}
func (uc *taskUcase) Store(ctx context.Context, p *ReqCreateTask) (Task, error) {
	panic("not implemented")
}
func (uc *taskUcase) Get(ctx context.Context, uid UserID, tid TaskID) (Task, error) {
	panic("not implemented")
}
func (uc *taskUcase) Patch(ctx context.Context, uid UserID, tid TaskID, p *ReqPatchTask) (Task, error) {
	panic("not implemented")
}
func (uc *taskUcase) Delete(ctx context.Context, uid UserID, tid TaskID) error {
	panic("not implemented")
}
