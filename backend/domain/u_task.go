package domain

type TaskID uint64

type TaskPatch struct {
	Title *string
	Done  *bool
}

type TaskListFilter struct {
	DoneEq *bool
}

type TaskUcase interface {
	Store(ctx Ctx, uid UserID, title string) (Task, error)
	List(Ctx, UserID, TaskListFilter) ([]Task, error)
	Get(Ctx, UserID, TaskID) (Task, error)
	Patch(Ctx, UserID, TaskID, TaskPatch) (Task, error)
	Delete(Ctx, UserID, TaskID) error
}
