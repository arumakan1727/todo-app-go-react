package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql/sqlcgen"
	"github.com/lib/pq"
)

func (r *repository) StoreUser(ctx context.Context, u *domain.User) error {
	now := r.clk.Now()
	id, err := r.q.InsertUser(ctx, r.db, sqlcgen.InsertUserParams{
		Email:       u.Email,
		Role:        u.Role,
		PasswdHash:  u.PasswdHash,
		DisplayName: u.DisplayName,
		CreatedAt:   now,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == PqErrCodeUniqueViolation {
				return fmt.Errorf("cannot create same email user(email='%s'): %w", u.Email, domain.ErrAlreadyExits)
			}
		}
		return err
	}
	u.ID = id
	u.CreatedAt = now
	return nil
}

func (r *repository) ListUsers(ctx context.Context) ([]domain.User, error) {
	got, err := r.q.ListUsers(ctx, r.db)

	resp := make([]domain.User, 0, len(got))
	for i := range got {
		u := &got[i]
		resp = append(resp, domain.User{
			ID:          u.ID,
			Role:        u.Role,
			Email:       u.Email,
			PasswdHash:  nil,
			DisplayName: u.DisplayName,
			CreatedAt:   u.CreatedAt.In(r.clk.Location()),
		})
	}
	return resp, err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	got, err := r.q.GetUserByEmail(ctx, r.db, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return got, fmt.Errorf("no such email user: %w", domain.ErrNotFound)
		}
		return got, err
	}
	got.CreatedAt = got.CreatedAt.In(r.clk.Location())
	return got, nil
}
