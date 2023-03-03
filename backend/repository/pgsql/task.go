package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/arumakan1727/todo-app-go-react/domain"
	. "github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/repository/pgsql/sqlcgen"
)

func wrapTaskError(err error, tid TaskID) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("cannot find task(id=%d): %w", tid, domain.ErrNotFound)
	}
	return err
}

func (r *repository) StoreTask(
	ctx context.Context, uid UserID, title string,
) (Task, error) {
	task, err := r.q.InsertTask(ctx, r.db, sqlcgen.InsertTaskParams{
		UserID:    uid,
		Title:     title,
		CreatedAt: r.clk.Now(),
	})
	task.ApplyTimezone(r.clk.Location())
	return task, err
}

func (r *repository) ListTasks(
	ctx context.Context, uid UserID, f TaskListFilter,
) ([]Task, error) {
	query := `select id, user_id, title, done, created_at from tasks where user_id=$1`

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
	for i := range tasks {
		tasks[i].ApplyTimezone(r.clk.Location())
	}
	return tasks, nil
}

func (r *repository) GetTask(
	ctx context.Context, uid UserID, tid TaskID,
) (Task, error) {
	task, err := r.q.GetTask(ctx, r.db, sqlcgen.GetTaskParams{
		ID:     tid,
		UserID: uid,
	})
	if err != nil {
		return task, wrapTaskError(err, tid)
	}
	task.ApplyTimezone(r.clk.Location())
	return task, nil
}

func (r *repository) PatchTask(
	ctx context.Context, uid UserID, tid TaskID, p TaskPatch,
) (Task, error) {
	query := `UPDATE tasks SET `
	qParts := []string{}
	args := []any{tid, uid}

	if title, ok := p.Title.Take(); ok {
		args = append(args, title)
		qParts = append(qParts, fmt.Sprintf("title=$%d", len(args)))
	}
	if done, ok := p.Done.Take(); ok {
		args = append(args, done)
		qParts = append(qParts, fmt.Sprintf("done=$%d", len(args)))
	}
	if len(qParts) == 0 {
		return Task{}, ErrEmptyPatch
	}

	query += strings.Join(qParts, ", ")
	query += ` where id=$1 AND user_id=$2 returning *;`

	var task Task
	err := r.db.GetContext(ctx, &task, query, args...)
	if err != nil {
		return task, wrapTaskError(err, tid)
	}
	task.ApplyTimezone(r.clk.Location())
	return task, nil
}

func (r *repository) DeleteTask(
	ctx context.Context, uid UserID, tid TaskID,
) (Task, error) {
	task, err := r.q.DeleteTask(ctx, r.db, sqlcgen.DeleteTaskParams{
		ID:     tid,
		UserID: uid,
	})
	if err != nil {
		return task, wrapTaskError(err, tid)
	}
	task.ApplyTimezone(r.clk.Location())
	return task, nil
}
