package pgsql

import (
	"context"

	"github.com/arumakan1727/todo-app-go-react/domain"
)

func (r *repository) StoreUser(ctx context.Context, u *domain.User) error {
	panic("TODO")
}

func (r *repository) ListUsers(ctx context.Context) ([]domain.User, error) {
	panic("TODO")
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	panic("TODO")
}
