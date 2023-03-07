package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql/sqlcgen"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dbInterface interface {
	sqlcgen.DBTX

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)

	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type repository struct {
	db_internal *sqlx.DB
	tx_internal *sqlx.Tx // 非トランザクション中は nil 。

	// db はトランザクション中は tx_internal を参照し、
	// そうでないときは db_internal を参照する。プロキシ的な役割。
	// BeginTx(), CommitTx(), RollbackTx() によって書き換わる。
	db dbInterface

	q   *sqlcgen.Queries
	clk clock.Clocker
}

func NewRepository(
	ctx context.Context, cfg *config.Config, clk clock.Clocker,
) (domain.Repository, error) {
	db, err := openDB(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgsql.NewRepository: failed to open db: %w", err)
	}

	dbx := sqlx.NewDb(db, "postgres")
	return &repository{
		db_internal: dbx,
		tx_internal: nil,

		db:  dbx,
		q:   sqlcgen.New(),
		clk: clk.In(time.UTC),
	}, nil
}

func (r *repository) Close() {
	if r == nil || r.db_internal == nil {
		return
	}
	_ = r.db_internal.Close()
}

// TruncateAll は TRUNCATE 文を用いてテーブルから全レコードを効率的に削除する。
// シーケンスもリセットする。
// 対象のテーブルが FOREIGN KEY によってテーブル A から参照されている場合は、A もクリアする。
func (r *repository) TruncateAll(ctx context.Context) error {
	tables := []string{
		"users",
		"tasks",
	}
	query := fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE;`, strings.Join(tables, ","))
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("(*pgsql.Repository).TruncateAll: failed to exec '%s': %#v", query, err)
	}
	return nil
}

func openDB(
	ctx context.Context, cfg *config.Config,
) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PgSQLURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
