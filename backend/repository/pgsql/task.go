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
) ([]*Task, error) {
	query := `select id, title, done, created_at where user_id=$1 `

	switch done {
	case domain.TaskStatusFilterAny:
		break // do nothing
	case domain.TaskStatusFilterDone:
		query += " AND done=true"
	case domain.TaskStatusFilterTodo:
		query += " AND done=false"
	}

	tasks := make([]*Task, 0, 5)
	rows, err := q.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	if err := sqlx.StructScan(rows, tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (q *Queries) PatchTask(
	ctx context.Context,
	taskID int32,
	userID int32,
	p *domain.ReqPatchTask,
) (*Task, error) {
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
	if len(args) <= 0 {
		return nil, ErrNoFieldsUpdate
	}

	query += strings.Join(qParts, ", ")
	query += `) where id=$1 AND user_id=$2
	returning (title, done, created_at);`

	row := q.db.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return nil, err
	}

	t := Task{
		ID:     taskID,
		UserID: userID,
	}
	row.Scan(&t.Title, &t.Done, &t.CreatedAt)
	return &t, nil
}
