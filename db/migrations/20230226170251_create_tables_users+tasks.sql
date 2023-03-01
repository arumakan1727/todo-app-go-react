-- migrate:up
create extension citext;

create table users (
  id bigint generated always as identity primary key,
  role          text        not null default 'user',
  email         citext      not null,
  passwd_hash   bytea       not null,
  display_name  text        not null,
  created_at    timestamp   not null default current_timestamp
);

alter table users add constraint "users_email_key" unique (email);

create table tasks (
  id bigint generated always as identity primary key,
  user_id       bigint      not null,
  title         text        not null,
  done          boolean     not null default false,
  created_at    timestamp   not null default current_timestamp
);

alter table tasks add constraint "tasks_user_id_fkey"
  foreign key (user_id) references users ("id")
  on update cascade
  on delete cascade;

-- migrate:down
drop table if exists users, tasks;
