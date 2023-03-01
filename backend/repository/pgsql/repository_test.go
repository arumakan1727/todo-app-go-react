package pgsql

import (
	"context"
	"fmt"
	"testing"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
)

func newRepositoryForTest(t *testing.T, ctx context.Context, clk clock.Clocker) domain.Repository {
	t.Helper()
	r, cleanup, err := NewRepository(ctx, config.ForTesting(), clk)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)
	return r
}

func clearTable(t *testing.T, ctx context.Context, r *repository, table string) {
	t.Helper()
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(`DELETE FROM "%s";`, table))
	if err != nil {
		t.Fatalf(`failed to exec 'DELETE FROM "%s"': %#v`, table, err)
	}
}
