package usecase

import (
	. "github.com/arumakan1727/todo-app-go-react/domain"
)

type taskUc struct {
	repo Repository
}

func NewTaskUsecase(repo Repository) TaskUcase {
	return &taskUc{
		repo: repo,
	}
}

func (uc *taskUc) Store(ctx Ctx, uid UserID, title string) (Task, error) {
	t, err := uc.repo.StoreTask(ctx, uid, title)
	return t, err
}

func (uc *taskUc) List(ctx Ctx, uid UserID, f TaskListFilter) ([]Task, error) {
	ts, err := uc.repo.ListTasks(ctx, uid, f)
	return ts, err
}

func (uc *taskUc) Get(ctx Ctx, uid UserID, tid TaskID) (Task, error) {
	t, err := uc.repo.GetTask(ctx, uid, tid)
	return t, err
}

func (uc *taskUc) Patch(ctx Ctx, uid UserID, tid TaskID, p TaskPatch) (Task, error) {
	t, err := uc.repo.PatchTask(ctx, uid, tid, p)
	return t, err
}

func (uc *taskUc) Delete(ctx Ctx, uid UserID, tid TaskID) (Task, error) {
	t, err := uc.repo.DeleteTask(ctx, uid, tid)
	return t, err
}
