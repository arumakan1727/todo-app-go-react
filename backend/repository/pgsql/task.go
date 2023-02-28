package pgsql

import (
	"context"
	"fmt"
	"strings"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/jmoiron/sqlx"
)

func (q *Queries) ListTasksFilterByStatus(
	ctx context.Context,
	userID int32,
	done domain.TaskStatusFilter,
) ([]domain.TaskEntity, error) {
	query := `select id, title, done, created_at where user_id=$1 `

	switch done {
	case domain.TaskStatusFilterAny:
		break // do nothing
	case domain.TaskStatusFilterDone:
		query += " AND done=true"
	case domain.TaskStatusFilterTodo:
		query += " AND done=false"
	}

	rows, err := q.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.TaskEntity, 0, 5)
	if err := sqlx.StructScan(rows, tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (q *Queries) PatchTask(
	ctx context.Context,
	taskID domain.TaskID,
	userID domain.UserID,
	p *domain.ReqPatchTask,
) (domain.TaskEntity, error) {
	query := `UPDATE tasks SET `
	qParts := []string{}
	args := []any{taskID, userID}

	if p.Title != nil {
		args = append(args, p.Title)
		qParts = append(qParts, fmt.Sprintf("title = $%d", len(args)))
	}
	if p.Done != nil {
		args = append(args, p.Done)
		qParts = append(qParts, fmt.Sprintf("done = $%d", len(args)))
	}
	if len(args) == 0 {
		return domain.TaskEntity{}, ErrNoFieldsUpdate
	}

	query += strings.Join(qParts, ", ")
	query += `) where id=$1 AND user_id=$2
	returning (title, done, created_at);`

	row := q.db.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return domain.TaskEntity{}, err
	}

	t := domain.TaskEntity{
		ID:     taskID,
		UserID: userID,
	}
	if err := row.Scan(&t.Title, &t.Done, &t.CreatedAt); err != nil {
		return domain.TaskEntity{}, err
	}
	return t, nil
}
