package domain

type Repository interface {
	StoreUser(Ctx, *User) error
	ListUsers(Ctx) ([]User, error)
	GetUserByEmail(Ctx, string) (User, error)

	StoreTask(Ctx, *Task) error
	ListTasksOfUser(Ctx, UserID) ([]Task, error)
	GetTask(Ctx, UserID, TaskID) (Task, error)
	PatchTask(Ctx, UserID, TaskPatch) (Task, error)
	DeleteTask(Ctx, UserID, TaskID) error
}
