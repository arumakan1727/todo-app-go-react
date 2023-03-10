-- name: GetUserByEmail :one
select * from users where email = $1 limit 1;

-- name: ListUsers :many
select
  id, role, email, display_name, created_at
from users;

-- name: InsertUser :one
insert into users (
  email, role, display_name, passwd_hash, created_at
) values ($1, $2, $3, $4, $5)
returning id;

-- name: GetTask :one
select * from tasks where id = $1 and user_id=$2 limit 1;

-- name: InsertTask :one
insert into tasks (
   user_id, title, created_at
) values ($1, $2, $3)
returning *;

-- name: DeleteTask :one
delete from tasks where id = $1 and user_id=$2 returning *;
