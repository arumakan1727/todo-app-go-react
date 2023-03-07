package domain

type Repository interface {
	Close()

	// BeginTx はトランザクションを開始する。
	// すでにトランザクション中の場合は何もしない (非エラー)。
	BeginTx(Ctx) error

	// CommitTx はトランザクションをコミットする。
	CommitTx(Ctx) error

	// RollbackTx はトランザクションをロールバックする。
	RollbackTx(Ctx) error

	StoreUser(Ctx, *User) error
	ListUsers(Ctx) ([]User, error)
	GetUserByEmail(Ctx, string) (User, error)

	StoreTask(ctx Ctx, uid UserID, title string) (Task, error)
	ListTasks(Ctx, UserID, TaskListFilter) ([]Task, error)
	GetTask(Ctx, UserID, TaskID) (Task, error)
	PatchTask(Ctx, UserID, TaskID, TaskPatch) (Task, error)
	DeleteTask(Ctx, UserID, TaskID) (Task, error)
}
