-- WARN: This file is only for sqlc.
-- 構造体生成に必要な情報しか記述していないので実際のDBで実行しないこと。

-- NOTE: ID の型名 (user_id 等) sqlc で Go の defined type を指定するため。
-- 通常の bigint や serial と区別するために固有の型名を指定してある。

create table users (
  id            user_id     primary key,
  role          text        not null,
  email         text        not null,
  passwd_hash   bytea       not null,
  display_name  text        not null,
  created_at    timestamptz not null
);

create table tasks (
  id            task_id     primary key,
  user_id       user_id     not null,
  title         text        not null,
  done          boolean     not null,
  created_at    timestamptz not null
);
