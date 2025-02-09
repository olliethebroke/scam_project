-- +goose Up
create table if not exists completed_tasks (
    id serial primary key,
    user_id bigint,
    task_id smallint
);

-- +goose Down
drop table completed_tasks;
