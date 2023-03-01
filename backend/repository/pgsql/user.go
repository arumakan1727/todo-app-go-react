package pgsql

import (
	"context"
	"errors"
	"fmt"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql/sqlcgen"
	"github.com/lib/pq"
)

func (r *repository) StoreUser(ctx context.Context, u *domain.User) error {
	id, err := r.q.InsertUser(ctx, r.db, sqlcgen.InsertUserParams{
		Email:       u.Email,
		DisplayName: u.DisplayName,
		PasswdHash:  u.PasswdHash,
		CreatedAt:   r.clk.Now(),
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == PqErrCodeUniqueViolation {
				return fmt.Errorf("cannot create same email user: %w", domain.ErrAlreadyExits)
			}
		}
		return err
	}
	u.ID = id
	return nil
}

func (r *repository) ListUsers(ctx context.Context) ([]domain.User, error) {
	panic("TODO")
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	panic("TODO")
}
