// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: query.sql

package pgsql

import (
	"context"
	"time"

	"github.com/arumakan1727/todo-app-go-react/domain"
)

const deleteTask = `-- name: DeleteTask :exec
delete from tasks where id = $1
`

func (q *Queries) DeleteTask(ctx context.Context, id domain.TaskID) error {
	_, err := q.db.ExecContext(ctx, deleteTask, id)
	return err
}

const getTask = `-- name: GetTask :one
select id, user_id, title, done, created_at from tasks where id = $1 limit 1
`

func (q *Queries) GetTask(ctx context.Context, id domain.TaskID) (domain.TaskEntity, error) {
	row := q.db.QueryRowContext(ctx, getTask, id)
	var i domain.TaskEntity
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Done,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, role, email, passwd_hash, display_name, created_at from users where email = $1 limit 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (domain.UserEntity, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i domain.UserEntity
	err := row.Scan(
		&i.ID,
		&i.Role,
		&i.Email,
		&i.PasswdHash,
		&i.DisplayName,
		&i.CreatedAt,
	)
	return i, err
}

const insertTask = `-- name: InsertTask :exec
insert into tasks (
   user_id, title, created_at
) values ($1, $2, $3)
`

type InsertTaskParams struct {
	UserID    domain.UserID `db:"user_id"`
	Title     string        `db:"title"`
	CreatedAt time.Time     `db:"created_at"`
}

func (q *Queries) InsertTask(ctx context.Context, arg InsertTaskParams) error {
	_, err := q.db.ExecContext(ctx, insertTask, arg.UserID, arg.Title, arg.CreatedAt)
	return err
}

const insertUser = `-- name: InsertUser :exec
insert into users (
  email, display_name, passwd_hash, created_at
) values ($1, $2, $3, $4)
`

type InsertUserParams struct {
	Email       string    `db:"email"`
	DisplayName string    `db:"display_name"`
	PasswdHash  []byte    `db:"passwd_hash"`
	CreatedAt   time.Time `db:"created_at"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) error {
	_, err := q.db.ExecContext(ctx, insertUser,
		arg.Email,
		arg.DisplayName,
		arg.PasswdHash,
		arg.CreatedAt,
	)
	return err
}

const listUsers = `-- name: ListUsers :many
select
  id, email, display_name, created_at
from users
`

type ListUsersRow struct {
	ID          domain.UserID `db:"id"`
	Email       string        `db:"email"`
	DisplayName string        `db:"display_name"`
	CreatedAt   time.Time     `db:"created_at"`
}

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.DisplayName,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
