package domain

type UserUsecase interface {
	List(ctx) (UserList, error)
	Store(ctx, *ReqCreateUser) (User, error)
}

type AuthTokenUsecase interface {
	Issue(ctx, *ReqCreateAuthToken) (AuthToken, error)
}

type TaskUsecase interface {
	List(ctx, ListTasksParams) (TaskList, error)
	Store(ctx, *ReqCreateTask) (Task, error)
	Get(ctx, UserID, TaskID) (Task, error)
	Patch(ctx, UserID, TaskID, *ReqPatchTask) (Task, error)
	Delete(ctx, UserID, TaskID) error
}
