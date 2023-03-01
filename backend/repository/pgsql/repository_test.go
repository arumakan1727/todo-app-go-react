package pgsql

import (
	"context"
	"testing"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
)

func newRepositoryForTest(t *testing.T, ctx context.Context, clk clock.Clocker, autoRollback bool) domain.Repository {
	t.Helper()
	r, cleanup, err := NewRepository(ctx, config.ForTesting(), clk)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)

	if !autoRollback {
		return r
	}

	if err := r.BeginTx(ctx); err != nil {
		t.Fatalf("newRepositoryForTest: %v", err)
	}
	t.Cleanup(func() {
		if err := r.RollbackTx(ctx); err != nil {
			t.Fatalf("Cleanup: newRepositoryForTest: %v", err)
		}
	})
	return r
}
