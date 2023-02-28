package pgsql

import (
	"context"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql/sqlcgen"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	dbx *sqlx.DB
	q   sqlcgen.Queries
	clk clock.Clocker
}

func NewRepository(ctx context.Context) domain.Repository {
	panic("TODO")
}
