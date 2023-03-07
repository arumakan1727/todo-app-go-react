package pgsql

import (
	"context"
	"testing"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
)

func newRepositoryForTest(t *testing.T, ctx context.Context, clk clock.Clocker) domain.Repository {
	t.Helper()
	r, err := NewRepository(ctx, config.ForTesting(), clk)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(r.Close)
	return r
}
