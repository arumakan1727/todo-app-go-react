package pgsql

import (
	"context"
	"fmt"
	"strings"

	"github.com/arumakan1727/todo-app-go-react/domain"
)

func (r *repository) ListTasks(
	ctx context.Context, userID domain.UserID, filterDoneEq *bool,
) ([]domain.Task, error) {
	query := `select id, title, done, created_at where user_id=$1 `

	if filterDoneEq != nil {
		if *filterDoneEq {
			query += " AND done=true"
		} else {
			query += " AND done=false"
		}
	}

	tasks := []domain.Task{}
	if err := r.dbx.SelectContext(ctx, &tasks, query, userID); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (q *repository) PatchTask(
	ctx context.Context,
	userID domain.UserID,
	taskID domain.TaskID,
	p domain.TaskPatch,
) (domain.Task, error) {
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
		return domain.Task{}, domain.ErrEmptyPatch
	}

	query += strings.Join(qParts, ", ")
	query += `) where id=$1 AND user_id=$2 returning *;`

	row := q.dbx.QueryRowxContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return domain.Task{}, err
	}

	var t domain.Task
	if err := row.StructScan(&t); err != nil {
		return domain.Task{}, err
	}
	return t, nil
}
