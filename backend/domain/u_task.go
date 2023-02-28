package domain

type TaskID uint64

type TaskPatch struct {
	Title *string
	Done  *bool
}

type TaskUcase interface {
	Store(ctx Ctx, uid UserID, title string) (Task, error)
	List(ctx Ctx, uid UserID, filterDoneEq *bool) ([]Task, error)
	Get(Ctx, UserID, TaskID) (Task, error)
	Patch(Ctx, UserID, TaskID, TaskPatch) (Task, error)
	Delete(Ctx, UserID, TaskID) error
}
