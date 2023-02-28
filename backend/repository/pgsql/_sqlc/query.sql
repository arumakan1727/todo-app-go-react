-- name: GetUserByEmail :one
select * from users where email = $1 limit 1;

-- name: ListUsers :many
select
  id, email, display_name, created_at
from users;

-- name: InsertUser :exec
insert into users (
  email, display_name, passwd_hash, created_at
) values ($1, $2, $3, $4);

-- name: GetTask :one
select * from tasks where id = $1 limit 1;

-- name: InsertTask :exec
insert into tasks (
   user_id, title, created_at
) values ($1, $2, $3);

-- name: DeleteTask :exec
delete from tasks where id = $1;
