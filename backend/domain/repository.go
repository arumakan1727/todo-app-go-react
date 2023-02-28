package domain

type Repository interface {
	StoreUser(ctx, *UserEntity) error
	ListUsers(ctx) ([]UserEntity, error)
	GetUserByEmail(ctx, Email) (UserEntity, error)

	StoreTask(ctx, *TaskEntity) error
	ListTasksOfUser(ctx, UserID) ([]TaskEntity, error)
	GetTask(ctx, UserID, TaskID) (TaskEntity, error)
	PatchTask(ctx, UserID, *ReqPatchTask) (TaskEntity, error)
	DeleteTask(ctx, UserID, TaskID) error
}
