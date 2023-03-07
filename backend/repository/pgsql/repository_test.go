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
	r, err := NewRepository(ctx, config.ForTesting(), clk)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(r.Close)
	return r
}

// clearTable は TRUNCATE 文を用いてテーブルから全レコードを効率的に削除する。
// シーケンスもリセットする。
// 対象のテーブルが FOREIGN KEY によってテーブル A から参照されている場合は、A もクリアする。
func clearTable(t *testing.T, ctx context.Context, r *repository, table string) {
	t.Helper()
	query := fmt.Sprintf(`TRUNCATE TABLE "%s" RESTART IDENTITY CASCADE;`, table)
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		t.Fatalf(`failed to exec '%s': %#v`, query, err)
	}
}
