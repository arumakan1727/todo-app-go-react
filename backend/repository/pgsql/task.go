package pgsql

import (
	"context"
	"fmt"
	"strings"

	. "github.com/arumakan1727/todo-app-go-react/domain"
)

func (r *repository) StoreTask(
	ctx context.Context, t *Task,
) error {
	panic("TODO")
}

func (r *repository) ListTasks(
	ctx context.Context, uid UserID, f TaskListFilter,
) ([]Task, error) {
	query := `select id, title, done, created_at where user_id=$1 `

	if doneEq, ok := f.DoneEq.Take(); ok {
		if doneEq {
			query += " AND done=true"
		} else {
			query += " AND done=false"
		}
	}

	tasks := []Task{}
	if err := r.db.SelectContext(ctx, &tasks, query, uid); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *repository) GetTask(
	ctx context.Context, uid UserID, tid TaskID,
) (Task, error) {
	panic("TODO")
}

func (r *repository) PatchTask(
	ctx context.Context, uid UserID, tid TaskID, p TaskPatch,
) (Task, error) {
	query := `UPDATE tasks SET `
	qParts := []string{}
	args := []any{tid, uid}

	if title, ok := p.Title.Take(); ok {
		args = append(args, title)
		qParts = append(qParts, fmt.Sprintf("title = $%d", len(args)))
	}
	if done, ok := p.Done.Take(); ok {
		args = append(args, done)
		qParts = append(qParts, fmt.Sprintf("done = $%d", len(args)))
	}
	if len(args) == 0 {
		return Task{}, ErrEmptyPatch
	}

	query += strings.Join(qParts, ", ")
	query += `) where id=$1 AND user_id=$2 returning *;`

	row := r.db.QueryRowxContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return Task{}, err
	}

	var t Task
	if err := row.StructScan(&t); err != nil {
		return Task{}, err
	}
	return t, nil
}

func (r *repository) DeleteTask(
	ctx context.Context, uid UserID, tid TaskID,
) error {
	panic("TODO")
}
