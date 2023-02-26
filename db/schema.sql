create extension citext; -- case insentive text

create table users (
  id            serial      primary key,
  role          text        not null default 'user',
  email         citext      not null,
  passwd_hash   text        not null,
  display_name  text        not null,
  created_at    timestamp   not null default current_timestamp
);

alter table users add constraint "users_email_key" UNIQUE (email);

create table tasks (
  id            serial      primary key,
  user_id       int         not null,
  title         text        not null,
  done          boolean     not null default false,
  created_at    timestamp   not null default current_timestamp
);

alter table tasks add constraint "tasks_user_id_fkey"
  foreign key (user_id) REFERENCES users ("id")
  on update cascade
  on delete cascade;
